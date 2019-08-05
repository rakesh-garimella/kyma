package serving

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/avast/retry-go"
	"github.com/knative/serving/pkg/apis/serving/v1beta1"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"

	serving_api "github.com/knative/serving/pkg/apis/serving/v1alpha1"
	serving "github.com/knative/serving/pkg/client/clientset/versioned/typed/serving/v1alpha1"
	"github.com/kyma-project/kyma/common/ingressgateway"
	corev1 "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ServingUpgradeTest struct {
	functionName, uuid string
	coreClient         kubernetes.Interface
	nSpace             string
	domain             string
	stop               <-chan struct{}
	servingClient      serving.ServiceInterface
}

// NewLambdaFunctionUpgradeTest returns new instance of the FunctionUpgradeTest
func NewServingUpgradeTest(servingCli serving.ServiceInterface, k8sCli kubernetes.Interface, domainName string) *ServingUpgradeTest {
	nSpace := strings.ToLower("ServingUpgradeTest")

	return &ServingUpgradeTest{
		servingClient: servingCli,
		coreClient:    k8sCli,
		nSpace:        nSpace,
		domain:        domainName,
	}
}

// CreateResources creates resources needed for e2e upgrade test
func (s *ServingUpgradeTest) CreateResources(stop <-chan struct{}, log logrus.FieldLogger, namespace string) error {
	log.Println("ServingUpgradeTest creating resources")
	s.nSpace = namespace
	s.stop = stop

	_, err := s.createServing()
	if err != nil {
		return err
	}

	// Ensure resources works
	err = s.TestResources(stop, log, namespace)
	if err != nil {
		return errors.Wrap(err, "first call to TestResources() failed.")
	}
	return nil
}

func (s *ServingUpgradeTest) createServing() (*serving_api.Service, error) {
	service, err := s.servingClient.Create(&serving_api.Service{
		ObjectMeta: meta.ObjectMeta{
			Name: "test-service",
		},
		Spec: serving_api.ServiceSpec{
			ConfigurationSpec: serving_api.ConfigurationSpec{
				Template: &serving_api.RevisionTemplateSpec{
					ObjectMeta: meta.ObjectMeta{
						Name:      "test-service",
						Namespace: s.nSpace,
					},
					Spec: serving_api.RevisionSpec{
						RevisionSpec: v1beta1.RevisionSpec{
							PodSpec: v1beta1.PodSpec{
								Containers: []corev1.Container{{
									Image: "gcr.io/knative-samples/helloworld-go",
									Env: []corev1.EnvVar{
										{
											Name:  "TARGET",
											Value: "knservice",
										},
									},
								},
								},
							},
						},
					},
				},
			},
			RouteSpec: serving_api.RouteSpec{
				Traffic: []serving_api.TrafficTarget{{
					TrafficTarget: v1beta1.TrafficTarget{
						Tag:          "rev-01",
						RevisionName: "test-service",
						Percent:      100,
					},
				},
				},
			},
		},
	})
	return service, err
}

func (s *ServingUpgradeTest) TestResources(stop <-chan struct{}, log logrus.FieldLogger, namespace string) error {
	testServiceURL := fmt.Sprintf("https://test-service.%s.%s", s.nSpace, s.domain)
	ingressClient, err := ingressgateway.FromEnv().Client()

	err = retry.Do(func() error {
		log.Printf("Calling: %s", testServiceURL)
		resp, err := ingressClient.Get(testServiceURL)
		if err != nil {
			return err
		}

		bytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		msg := strings.TrimSpace(string(bytes))
		expectedMsg := fmt.Sprintf("Hello %s!", "knservice")
		log.Printf("Received %v: '%s'", resp.StatusCode, msg)

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("unexpected status code: %v", resp.StatusCode)
		}
		if msg != expectedMsg {
			return fmt.Errorf("unexpected response: '%s'", msg)
		}

		return nil
	}, retry.OnRetry(func(n uint, err error) {
		log.Printf("[%v] try failed: %s", n, err)
	}), retry.Attempts(20),
	)
	return err

}
