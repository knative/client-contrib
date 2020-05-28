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

package migrate

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	apiv1 "k8s.io/api/core/v1"
	api_errors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	clientset "k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc" // from https://github.com/kubernetes/client-go/issues/345
	"k8s.io/client-go/tools/clientcmd"
	"knative.dev/client-contrib/plugins/migration/pkg/command"
	serving_v1_api "knative.dev/serving/pkg/apis/serving/v1"
	serving_v1_client "knative.dev/serving/pkg/client/clientset/versioned/typed/serving/v1"
)

type migrateCmdFlags struct {
	Namespace             string
	KubeConfig            string
	DestinationKubeConfig string
	DestinationNamespace  string
	Force                 bool
	Delete                bool
}

var MaxUpdateRetries = 3
var migrateFlags migrateCmdFlags

// migrateCmd represents the migrate command
func NewMigrateCommand() *cobra.Command {
	var migrateCmd = &cobra.Command{
		Use:   "migrate",
		Short: "Migrate Knative services from source cluster to destination cluster",
		Example: `
  # Migrate Knative services from source cluster to destination cluster by export KUBECONFIG and KUBECONFIG_DESTINATION as environment variables
  kn migrate --namespace default --destination-namespace default
  # Migrate Knative services from source cluster to destination cluster by set kubeconfig as parameters
  kn migrate --namespace default --destination-namespace default --kubeconfig $HOME/.kube/config/source-cluster-config.yml --destination-kubeconfig $HOME/.kube/config/destination-cluster-config.yml
  # Migrate Knative services from source cluster to destination cluster and force replace the service if exists in destination cluster
  kn migrate --namespace default --destination-namespace default --force
  # Migrate Knative services from source cluster to destination cluster and delete the service in source cluster
  kn migrate --namespace default --destination-namespace default --force --delete`,

		Run: func(cmd *cobra.Command, args []string) {
			kubeconfigS := migrateFlags.KubeConfig
			if kubeconfigS == "" {
				kubeconfigS = os.Getenv("KUBECONFIG")
			}
			if kubeconfigS == "" {
				fmt.Printf("cannot get source cluster kube config, please use --kubeconfig or export environment variable KUBECONFIG to set\n")
				os.Exit(1)
			}

			kubeconfigD := migrateFlags.DestinationKubeConfig
			if kubeconfigD == "" {
				kubeconfigD = os.Getenv("KUBECONFIG_DESTINATION")
			}
			if kubeconfigD == "" {
				fmt.Printf("cannot get destination cluster kube config, please use --destination-kubeconfig or export environment variable KUBECONFIG_DESTINATION to set\n")
				os.Exit(1)
			}

			namespaceS := migrateFlags.Namespace
			if namespaceS == "" {
				fmt.Printf("cannot get source cluster namespace, please use --namespace to set\n")
				os.Exit(1)
			}

			namespaceD := migrateFlags.DestinationNamespace
			if namespaceD == "" {
				fmt.Printf("cannot get destination cluster namespace, please use --destination-namespace to set\n")
				os.Exit(1)
			}

			// For source
			_, migrationClientS, err := getClients(kubeconfigS, namespaceS)
			if err != nil {
				fmt.Printf(err.Error())
				os.Exit(1)
			}
			err = migrationClientS.PrintServiceWithRevisions("source")
			if err != nil {
				fmt.Printf(err.Error())
				os.Exit(1)
			}

			// For destination
			clientSetD, migrationClientD, err := getClients(kubeconfigD, namespaceD)
			if err != nil {
				fmt.Printf(err.Error())
				os.Exit(1)
			}

			fmt.Println(color.GreenString("[Before migration in destination cluster]"))
			err = migrationClientD.PrintServiceWithRevisions("destination")
			if err != nil {
				fmt.Printf(err.Error())
				os.Exit(1)
			}

			fmt.Println("\nNow migrate all Knative service resources")
			fmt.Println("From the source", color.BlueString(namespaceS), "namespace of cluster", color.CyanString(kubeconfigS))
			fmt.Println("To the destination", color.BlueString(namespaceD), "namespace of cluster", color.CyanString(kubeconfigD))

			err = getOrCreateNamespace(clientSetD, namespaceD)
			if err != nil {
				fmt.Printf(err.Error())
				os.Exit(1)
			}

			servicesS, err := migrationClientS.ListService()
			if err != nil {
				fmt.Printf(err.Error())
				os.Exit(1)
			}
			for i := 0; i < len(servicesS.Items); i++ {
				serviceS := servicesS.Items[i]
				err = createService(migrationClientD, serviceS, migrateFlags.Force)
				if err != nil {
					fmt.Printf(err.Error())
					os.Exit(1)
				}
				fmt.Println("Migrated service", color.CyanString(serviceS.Name), "Successfully")

				serviceD, err := migrationClientD.GetService(serviceS.Name)
				if err != nil {
					fmt.Printf(err.Error())
					os.Exit(1)
				}

				config, err := migrationClientD.GetConfig(serviceD.Name)
				if err != nil {
					fmt.Printf(err.Error())
					os.Exit(1)
				}
				configUuid := config.UID

				revisionsS, err := migrationClientS.ListRevisionByService(serviceS.Name)
				if err != nil {
					fmt.Printf(err.Error())
					os.Exit(1)
				}
				for i := 0; i < len(revisionsS.Items); i++ {
					revisionS := revisionsS.Items[i]
					err = migrateRevision(migrationClientD, revisionS, serviceS, configUuid)
					if err != nil {
						fmt.Printf(err.Error())
						os.Exit(1)
					}
				}
				fmt.Println("")
			}

			fmt.Println(color.GreenString("[After migration in destination cluster]"))
			err = migrationClientD.PrintServiceWithRevisions("destination")
			if err != nil {
				fmt.Printf(err.Error())
				os.Exit(1)
			}

			err = deleteAllServices(migrationClientS, migrateFlags.Delete)
			if err != nil {
				fmt.Printf(err.Error())
				os.Exit(1)
			}
		},
	}

	migrateCmd.Flags().StringVarP(&migrateFlags.Namespace, "namespace", "n", "", "The namespace of the source Knative resources")
	migrateCmd.Flags().StringVar(&migrateFlags.KubeConfig, "kubeconfig", "", "The kubeconfig of the Knative resources (default is KUBECONFIG from environment variable)")

	migrateCmd.Flags().StringVar(&migrateFlags.DestinationKubeConfig, "destination-kubeconfig", "", "The kubeconfig of the destination Knative resources (default is KUBECONFIG_DESTINATION from environment variable)")
	migrateCmd.Flags().StringVar(&migrateFlags.DestinationNamespace, "destination-namespace", "", "The namespace of the destination Knative resources")

	migrateCmd.Flags().BoolVar(&migrateFlags.Force, "force", false, "Migrate service forcefully, replaces existing service if any.")
	migrateCmd.Flags().BoolVar(&migrateFlags.Delete, "delete", false, "Delete all Knative resources after kn-migration from source cluster")
	return migrateCmd
}

func getClients(kubeConfig, namespace string) (*kubernetes.Clientset, command.MigrationClient, error) {
	cfg, err := clientcmd.BuildConfigFromFlags("", kubeConfig)
	if err != nil {
		return nil, nil, err
	}
	clientSet, err := clientset.NewForConfig(cfg)
	if err != nil {
		return nil, nil, err
	}
	servingClient, err := serving_v1_client.NewForConfig(cfg)
	if err != nil {
		return nil, nil, err
	}
	migrationClient := command.NewMigrationClient(servingClient, namespace)
	return clientSet, migrationClient, nil
}

func getOrCreateNamespace(clientSet *kubernetes.Clientset, namespace string) error {
	namespaceExists := true
	_, err := clientSet.CoreV1().Namespaces().Get(namespace, metav1.GetOptions{})
	if api_errors.IsNotFound(err) {
		namespaceExists = false
	}
	if err != nil {
		return err
	}

	if !namespaceExists {
		fmt.Println("Create namespace", color.BlueString(migrateFlags.Namespace), "in destination cluster")
		nsSpec := &apiv1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: namespace}}
		_, err := clientSet.CoreV1().Namespaces().Create(nsSpec)
		if err != nil {
			return err
		}
	} else {
		fmt.Println("Namespace", migrateFlags.Namespace, "already exists in destination cluster")
	}
	return nil
}

func createService(migrationClient command.MigrationClient, service serving_v1_api.Service, force bool) error {
	serviceExists, err := migrationClient.ServiceExists(service.Name)
	if err != nil {
		return err
	}

	if serviceExists {
		if !force {
			return fmt.Errorf("cannot migrate service %s in namespace because the service already exists and no --force option was given", service.Name)
		}
		fmt.Println("Deleting service", color.CyanString(service.Name), "from the destination cluster and recreate as replacement")
		migrationClient.DeleteService(service.Name)
		if err != nil {
			return err
		}
	}
	_, err = migrationClient.CreateService(&service)
	if err != nil {
		return err
	}
	return nil
}

func deleteAllServices(migrationClient command.MigrationClient, delete bool) error {
	if !delete {
		fmt.Println("Migrate without --delete option, skip deleting Knative resource in source cluster")
	} else {
		fmt.Println("Migrate with --delete option, deleting all Knative resource in source cluster")
		services, err := migrationClient.ListService()
		if err != nil {
			return err
		}
		for i := 0; i < len(services.Items); i++ {
			service := services.Items[i]
			err = migrationClient.DeleteService(service.Name)
			if err != nil {
				return err
			}
			fmt.Println("Deleted service", service.Name, "in source cluster")
		}
	}
	return nil
}

func migrateRevision(migrationClient command.MigrationClient, revisionS serving_v1_api.Revision, serviceS serving_v1_api.Service, configUuid types.UID) error {
	if revisionS.Name != serviceS.Status.LatestReadyRevisionName {
		_, err := migrationClient.CreateRevision(&revisionS, configUuid)
		if err != nil {
			return err
		}
		fmt.Println("Migrated revision", color.CyanString(revisionS.Name), "successfully")
	} else {
		retries := 0
		for {
			revision, err := migrationClient.GetRevision(revisionS.Name)
			if err != nil {
				return err
			}
			sourceRevisionGeneration := revisionS.ObjectMeta.Labels["serving.knative.dev/configurationGeneration"]
			revision.ObjectMeta.Labels["serving.knative.dev/configurationGeneration"] = sourceRevisionGeneration
			err = migrationClient.UpdateRevision(revision)
			if err != nil {
				// Retry to update when a resource version conflict exists
				if api_errors.IsConflict(err) && retries < MaxUpdateRetries {
					retries++
					continue
				}
				return err
			}
			fmt.Println("Replace revision", color.CyanString(revisionS.Name), "to generation", sourceRevisionGeneration, "successfully")
			break
		}
	}
	return nil
}
