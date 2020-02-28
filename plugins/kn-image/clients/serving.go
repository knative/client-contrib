package clients

import (
	v1 "k8s.io/api/core/v1"
	api_errors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	serving_v1_api "knative.dev/serving/pkg/apis/serving/v1"
	serving_v1_client "knative.dev/serving/pkg/client/clientset/versioned/typed/serving/v1"
)

var MaxUpdateRetries = 3

type ServingClient interface {
	// Create service struct from provided options
	ConstructService(name, image, serviceAccount, namespace string) *serving_v1_api.Service

	// Check if service exists
	ServiceExists(name string) (bool, *serving_v1_api.Service, error)

	// Get a service by name
	GetService(name string) (*serving_v1_api.Service, error)

	// Update a service
	UpdateService(service *serving_v1_api.Service, image, serviceAccount string) error

	// Create a service
	CreateService(service *serving_v1_api.Service) error
}

type servingClient struct {
	client    serving_v1_client.ServingV1Interface
	namespace string
}

// Create a new client facade for the provided cl.namespace
func NewServingClient(client serving_v1_client.ServingV1Interface, namespace string) ServingClient {
	return &servingClient{
		client:    client,
		namespace: namespace,
	}
}

// Create service struct from provided options
func (cl *servingClient) ConstructService(name, image, serviceAccount, namespace string) *serving_v1_api.Service {

	service := serving_v1_api.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: serving_v1_api.ServiceSpec{
			ConfigurationSpec:    serving_v1_api.ConfigurationSpec{
				Template: serving_v1_api.RevisionTemplateSpec{
					Spec: serving_v1_api.RevisionSpec{
						PodSpec: v1.PodSpec{
							ServiceAccountName: serviceAccount,
							Containers: []v1.Container{
								{
									Image: image,
								},
							},
						},
					},
				},
			},
		},
	}

	return &service
}

// Check if service exists
func (cl *servingClient) ServiceExists(name string) (bool, *serving_v1_api.Service, error) {
	service, err := cl.client.Services(cl.namespace).Get(name, metav1.GetOptions{})
	if api_errors.IsNotFound(err) {
		return false, nil, nil
	}
	if err != nil {
		return false, nil, err
	}
	return true, service, nil
}

// Get a service by name
func (cl *servingClient) GetService(name string) (*serving_v1_api.Service, error) {
	service, err := cl.client.Services(cl.namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return service, nil
}

// Update a service
func (cl *servingClient) UpdateService(service *serving_v1_api.Service, image, serviceAccount string) error {
	var retries = 0
	for {
		service.Spec = serving_v1_api.ServiceSpec{
			ConfigurationSpec:    serving_v1_api.ConfigurationSpec{
				Template: serving_v1_api.RevisionTemplateSpec{
					Spec: serving_v1_api.RevisionSpec{
						PodSpec: v1.PodSpec{
							ServiceAccountName: serviceAccount,
							Containers: []v1.Container{
								{
									Image: image,
								},
							},
						},
					},
				},
			},
		}
		_, err := cl.client.Services(cl.namespace).Update(service)
		if err != nil {
			// Retry to update when a resource version conflict exists
			if api_errors.IsConflict(err) && retries < MaxUpdateRetries {
				retries++
				continue
			}
			return err
		}
		return nil
	}
}

// Create a service
func (cl *servingClient) CreateService(service *serving_v1_api.Service) error {
	_, err := cl.client.Services(cl.namespace).Create(service)
	if err != nil {
		return err
	}
	return nil
}

// Delete a service by name
func (cl *servingClient) DeleteService(name string) error {
	err := cl.client.Services(cl.namespace).Delete(name, &metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	return nil
}
