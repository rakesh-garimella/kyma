package backupe2e

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/avast/retry-go"
	serving_api "github.com/knative/serving/pkg/apis/serving/v1alpha1"
	"github.com/knative/serving/pkg/apis/serving/v1beta1"
	serving "github.com/knative/serving/pkg/client/clientset/versioned/typed/serving/v1alpha1"
	"github.com/kyma-project/kyma/common/ingressgateway"
	"github.com/kyma-project/kyma/tests/end-to-end/backup-restore-test/utils/config"
	corev1 "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type servingTest struct {
	servClient serving.ServiceInterface
}

func NewServingTest() (*servingTest, error) {
	restConfig, err := config.NewRestClientConfig()
	if err != nil {
		return &servingTest{}, err
	}

	serv := serving.NewForConfigOrDie(restConfig).Services("knative-serving")
	if err != nil {
		return &servingTest{}, err
	}

	return &servingTest{
		servClient: serv,
	}, nil
}

func (s *servingTest) createService(namespace, target, name string) (*serving_api.Service, error) {
	service, err := s.servClient.Create(&serving_api.Service{
		ObjectMeta: meta.ObjectMeta{
			Name: "test-service",
		},
		Spec: serving_api.ServiceSpec{
			ConfigurationSpec: serving_api.ConfigurationSpec{
				Template: &serving_api.RevisionTemplateSpec{
					ObjectMeta: meta.ObjectMeta{
						Name:      name,
						Namespace: namespace,
					},
					Spec: serving_api.RevisionSpec{
						RevisionSpec: v1beta1.RevisionSpec{
							PodSpec: v1beta1.PodSpec{
								Containers: []corev1.Container{{
									Image: "gcr.io/knative-samples/helloworld-go",
									Env: []corev1.EnvVar{
										{
											Name:  "TARGET",
											Value: target,
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
						RevisionName: name,
						Percent:      100,
					},
				},
				},
			},
		},
	})
	return service, err
}

func (s *servingTest) updateService(namespace, target, name string) (*serving_api.Service, error) {
	service, err := s.servClient.Get(name, meta.GetOptions{})
	if err != nil {
		return nil, err
	}

	trafficTarget := serving_api.TrafficTarget{
		TrafficTarget: v1beta1.TrafficTarget{
			Tag:          "rev-02",
			RevisionName: name,
			Percent:      100,
		},
	}
	service.Spec.RouteSpec.Traffic = append(service.Spec.RouteSpec.Traffic, trafficTarget)
	serviceUpdate, err := s.servClient.Update(service)
	if err != nil {
		log.Printf("Unable to update the service: %v", err)
		return nil, err
	}
	return serviceUpdate, err
}

func (s *servingTest) checkServingStatus(namespace, target, rev string) error {
	domainName := s.MustGetenv("DOMAIN")
	testServiceURL := ""
	if rev != "" {
		testServiceURL = fmt.Sprintf("https://%s.test-service.knative-serving.%s", rev, domainName)
	} else {
		testServiceURL = fmt.Sprintf("https://test-service.knative-serving.%s", domainName)
	}

	ingressClient, err := ingressgateway.FromEnv().Client()
	if err != nil {
		log.Fatalf("Unexpected error when creating ingressgateway client: %s", err)
	}

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
		expectedMsg := fmt.Sprintf("Hello %s!", target)
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

func (s *servingTest) MustGetenv(name string) string {
	env := os.Getenv(name)
	if env == "" {
		log.Fatalf("Missing '%s' variable", name)
	}
	return env
}

func (s *servingTest) CreateResources(namespace string) {
	_, err := s.createService(namespace, "Rev-01", "rev-01")
	So(err, ShouldBeNil)
	_, err = s.updateService(namespace, "Rev-02", "rev-02")
	So(err, ShouldBeNil)
}

func (s *servingTest) TestResources(namespace string) {
	err := s.checkServingStatus(namespace, "Rev-02", "")
	So(err, ShouldBeNil)
	err = s.checkServingStatus(namespace, "Rev-01", "rev-01")
	So(err, ShouldBeNil)
}
