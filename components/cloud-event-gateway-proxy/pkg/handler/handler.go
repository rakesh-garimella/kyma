package handler

import (
	"context"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/cloudevents/sdk-go/v2/binding"
	cev2client "github.com/cloudevents/sdk-go/v2/client"
	cev2event "github.com/cloudevents/sdk-go/v2/event"
	cehttp "github.com/cloudevents/sdk-go/v2/protocol/http"

	cloudevents "github.com/kyma-project/kyma/components/cloud-event-gateway-proxy/pkg/cloudevents"
	"github.com/kyma-project/kyma/components/cloud-event-gateway-proxy/pkg/ems"
	"github.com/kyma-project/kyma/components/cloud-event-gateway-proxy/pkg/health"
	"github.com/kyma-project/kyma/components/cloud-event-gateway-proxy/pkg/receiver"
	"github.com/kyma-project/kyma/components/cloud-event-gateway-proxy/pkg/sender"
)

const (
	// noDuration signals that the dispatch step hasn't started
	noDuration = -1
)

var (
	additionalHeaders = http.Header{
		"qos":    []string{string(ems.QosAtLeastOnce)},
		"Accept": []string{"application/json"},
	}
)

type Handler struct {
	// Receiver receives incoming HTTP requests
	Receiver *receiver.HttpMessageReceiver
	// Sender sends requests to the broker
	Sender *sender.HttpMessageSender
	// Defaulter sets default values to incoming events
	Defaulter cev2client.EventDefaulter
	// RequestTimeout timeout for outgoing requests
	RequestTimeout time.Duration
	// Logger default logger
	Logger *logrus.Logger
}

func NewHandler(receiver *receiver.HttpMessageReceiver, sender *sender.HttpMessageSender, requestTimeout time.Duration, logger *logrus.Logger) *Handler {
	return &Handler{Receiver: receiver, Sender: sender, RequestTimeout: requestTimeout, Logger: logger}
}

func (h *Handler) Start(ctx context.Context) error {
	return h.Receiver.StartListen(ctx, health.CheckHealth(h))
}

func (h *Handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	// validate request method
	if request.Method != http.MethodPost {
		h.Logger.Warnf("Unexpected request method: %s", request.Method)
		h.writeResponse(writer, http.StatusMethodNotAllowed, nil)
		return
	}

	// validate request URI
	if request.RequestURI != "/publish" {
		h.writeResponse(writer, http.StatusNotFound, nil)
		return
	}

	ctx, cancel := context.WithTimeout(request.Context(), h.RequestTimeout)
	defer cancel()

	message := cehttp.NewMessageFromHttpRequest(request)
	defer func() { _ = message.Finish(nil) }()

	event, err := binding.ToEvent(ctx, message)
	if err != nil {
		h.Logger.Warnf("Failed to extract event from request with error: %s", err)
		h.writeResponse(writer, http.StatusBadRequest, []byte(err.Error()))
		return
	}

	if err := event.Validate(); err != nil {
		h.Logger.Warnf("Request is invalid as per CE spec with error: %s", err)
		h.writeResponse(writer, http.StatusBadRequest, []byte(err.Error()))
		return
	}

	h.receive(ctx, event)
	statusCode, dispatchTime, respBody := h.send(ctx, event)
	h.writeResponse(writer, statusCode, respBody)

	h.Logger.Infof("Event dispatch statusCode:[%d] duration:[%s] responseBody:[%s]", statusCode, dispatchTime, respBody)
}

func (h *Handler) writeResponse(writer http.ResponseWriter, statusCode int, respBody []byte) {
	writer.WriteHeader(statusCode)

	if respBody == nil {
		return
	}
	if _, err := writer.Write(respBody); err != nil {
		h.Logger.Errorf("Failed to write response body with error: %s", err)
	}
}

func (h *Handler) receive(ctx context.Context, event *cev2event.Event) {
	if h.Defaulter != nil {
		newEvent := h.Defaulter(ctx, *event)
		event = &newEvent
	}

	h.Logger.Infof("Event received id: %s", event.ID())
}

func (h *Handler) send(ctx context.Context, event *cev2event.Event) (int, time.Duration, []byte) {
	request, err := h.Sender.NewRequestWithTarget(ctx, h.Sender.Target)
	if err != nil {
		h.Logger.Errorf("Failed to prepare a cloudevent request with error: %s", err)
		return http.StatusInternalServerError, noDuration, []byte{}
	}

	message := binding.ToMessage(event)
	defer func() { _ = message.Finish(nil) }()

	err = cloudevents.WriteRequestWithHeaders(ctx, message, request, additionalHeaders)
	if err != nil {
		h.Logger.Errorf("Failed to add additional headers to the request with error: %s", err)
		return http.StatusInternalServerError, noDuration, []byte{}
	}

	resp, dispatchTime, err := h.sendAndRecordDispatchTime(request)
	if err != nil {
		h.Logger.Errorf("Failed to send event and record dispatch time with error: %s", err)
		return http.StatusInternalServerError, dispatchTime, []byte{}
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		h.Logger.Errorf("Failed to read response body with error: %s", err)
		return http.StatusInternalServerError, dispatchTime, []byte{}
	}

	return resp.StatusCode, dispatchTime, body
}

// sendAndRecordDispatchTime sends a CloudEvent and records the time taken while sending.
func (h *Handler) sendAndRecordDispatchTime(request *http.Request) (*http.Response, time.Duration, error) {
	start := time.Now()
	resp, err := h.Sender.Send(request)
	dispatchTime := time.Since(start)
	return resp, dispatchTime, err
}
