// Copyright © 2020 The Knative Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package command

import (
	"fmt"

	"github.com/fatih/color"
	api_errors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	api_serving "knative.dev/serving/pkg/apis/serving"
	serving_v1_api "knative.dev/serving/pkg/apis/serving/v1"
	serving_v1_client "knative.dev/serving/pkg/client/clientset/versioned/typed/serving/v1"
)

type MigrationClient interface {
	// Create service struct from provided options
	ConstructService(originalservice serving_v1_api.Service) *serving_v1_api.Service

	// Create revision struct from provided options
	BuildRevision(originalrevision serving_v1_api.Revision, config_uuid types.UID) *serving_v1_api.Revision

	// Check if service exists
	ServiceExists(name string) (bool, error)

	// Get a config by service name
	GetConfig(name string) (*serving_v1_api.Configuration, error)

	// Get a service by name
	GetService(name string) (*serving_v1_api.Service, error)

	// Get service list
	ListService() (*serving_v1_api.ServiceList, error)

	// Create a service
	CreateService(service *serving_v1_api.Service) (*serving_v1_api.Service, error)

	// Delete a service by name
	DeleteService(name string) error

	// Get a revision by service name
	GetRevision(name string) (*serving_v1_api.Revision, error)

	// Create a revision
	CreateRevision(revision *serving_v1_api.Revision, config_uuid types.UID) (*serving_v1_api.Revision, error)

	// Update the given revision
	UpdateRevision(revision *serving_v1_api.Revision) error

	// Get revision list by service
	ListRevisionByService(name string) (*serving_v1_api.RevisionList, error)

	// Get service list with revisions
	PrintServiceWithRevisions(clustername string) error
}

type migrationClient struct {
	client    serving_v1_client.ServingV1Interface
	namespace string
}

// NewMigrationClient creates a new client facade for the provided cl.namespace
func NewMigrationClient(client serving_v1_client.ServingV1Interface, namespace string) MigrationClient {
	return &migrationClient{
		client:    client,
		namespace: namespace,
	}
}

func (mc *migrationClient) ConstructService(originalservice serving_v1_api.Service) *serving_v1_api.Service {

	service := serving_v1_api.Service{
		ObjectMeta: originalservice.ObjectMeta,
	}

	service.ObjectMeta.Namespace = mc.namespace

	service.Spec = originalservice.Spec
	service.Spec.Template.ObjectMeta.Name = originalservice.Status.LatestCreatedRevisionName
	service.ObjectMeta.ResourceVersion = ""

	return &service
}

func (mc *migrationClient) BuildRevision(originalrevision serving_v1_api.Revision, config_uuid types.UID) *serving_v1_api.Revision {
	revision := serving_v1_api.Revision{
		ObjectMeta: originalrevision.ObjectMeta,
	}

	revision.ObjectMeta.Namespace = mc.namespace
	revision.ObjectMeta.ResourceVersion = ""
	revision.ObjectMeta.OwnerReferences[0].UID = config_uuid
	revision.ObjectMeta.Labels["serving.knative.dev/configurationGeneration"] = originalrevision.ObjectMeta.Labels["serving.knative.dev/configurationGeneration"]
	revision.Spec = originalrevision.Spec

	return &revision
}

func (mc *migrationClient) ServiceExists(name string) (bool, error) {
	_, err := mc.client.Services(mc.namespace).Get(name, metav1.GetOptions{})
	if api_errors.IsNotFound(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (mc *migrationClient) GetConfig(name string) (*serving_v1_api.Configuration, error) {
	config, err := mc.client.Configurations(mc.namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return config, nil
}

func (mc *migrationClient) GetService(name string) (*serving_v1_api.Service, error) {
	service, err := mc.client.Services(mc.namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return service, nil
}

func (mc *migrationClient) ListService() (*serving_v1_api.ServiceList, error) {
	servicelist, err := mc.client.Services(mc.namespace).List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return servicelist, nil
}

func (mc *migrationClient) CreateService(service *serving_v1_api.Service) (*serving_v1_api.Service, error) {
	newserivce := mc.ConstructService(*service)
	service, err := mc.client.Services(mc.namespace).Create(newserivce)
	if err != nil {
		return nil, err
	}
	return service, nil
}

func (mc *migrationClient) DeleteService(name string) error {
	err := mc.client.Services(mc.namespace).Delete(name, &metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (mc *migrationClient) GetRevision(name string) (*serving_v1_api.Revision, error) {
	revision, err := mc.client.Revisions(mc.namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return revision, nil
}

func (mc *migrationClient) CreateRevision(revision *serving_v1_api.Revision, config_uuid types.UID) (*serving_v1_api.Revision, error) {
	newrevision := mc.BuildRevision(*revision, config_uuid)
	revision, err := mc.client.Revisions(mc.namespace).Create(newrevision)
	if err != nil {
		return nil, err
	}
	return revision, nil
}

func (mc *migrationClient) UpdateRevision(revision *serving_v1_api.Revision) error {
	_, err := mc.client.Revisions(mc.namespace).Update(revision)
	if err != nil {
		return err
	}
	return nil
}

func (mc *migrationClient) ListRevisionByService(name string) (*serving_v1_api.RevisionList, error) {
	revisions, err := mc.client.Revisions(mc.namespace).List(metav1.ListOptions{LabelSelector: api_serving.ServiceLabelKey + "=" + name})
	if err != nil {
		return nil, err
	}
	return revisions, nil
}

func (mc *migrationClient) PrintServiceWithRevisions(clustername string) error {
	services, err := mc.ListService()
	if err != nil {
		return err
	}

	fmt.Println("There are", color.CyanString("%v", len(services.Items)), "service(s) in", clustername, color.BlueString(mc.namespace), "namespace")
	for i := 0; i < len(services.Items); i++ {
		service := services.Items[i]
		color.Cyan("%-25s%-30s%-20s\n", "Name", "Current Revision", "Ready")
		fmt.Printf("%-25s%-30s%-20s\n", service.Name, service.Status.LatestReadyRevisionName, fmt.Sprint(service.Status.IsReady()))

		revisions_s, err := mc.ListRevisionByService(service.Name)
		if err != nil {
			return err
		}
		for i := 0; i < len(revisions_s.Items); i++ {
			revision_s := revisions_s.Items[i]
			fmt.Println("  |- Revision", revision_s.Name, "( Generation: "+fmt.Sprint(revision_s.Labels["serving.knative.dev/configurationGeneration"]), ", Ready:", revision_s.Status.IsReady(), ")")
		}
		fmt.Println("")
	}
	return nil
}
