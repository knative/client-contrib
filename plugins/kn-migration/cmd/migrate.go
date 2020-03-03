// Copyright © 2019 The Knative Authors
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

package cmd

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	apiv1 "k8s.io/api/core/v1"
	api_errors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientset "k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc" // from https://github.com/kubernetes/client-go/issues/345
	"k8s.io/client-go/tools/clientcmd"
	serving_v1_client "knative.dev/serving/pkg/client/clientset/versioned/typed/serving/v1"
	"os"
)

// migrateCmd represents the migrate command
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
		namespaceS := ""
		namespaceD := ""
		if cmd.Flag("namespace").Value.String() == "" {
			fmt.Printf("cannot get source cluster namespace, please use --namespace to set\n")
			os.Exit(1)
		} else {
			namespaceS = cmd.Flag("namespace").Value.String()
		}

		if cmd.Flag("destination-kubeconfig").Value.String() == "" {
			fmt.Printf("cannot get destination cluster namespace, please use --destination-namespace to set\n")
			os.Exit(1)
		} else {
			namespaceD = cmd.Flag("destination-namespace").Value.String()
		}

		kubeconfigS := cmd.Flag("kubeconfig").Value.String()
		if kubeconfigS == "" {
			kubeconfigS = os.Getenv("KUBECONFIG")
		}
		if kubeconfigS == "" {
			fmt.Printf("cannot get source cluster kube config, please use --kubeconfig or export environment variable KUBECONFIG to set\n")
			os.Exit(1)
		}

		kubeconfigD := cmd.Flag("destination-kubeconfig").Value.String()
		if kubeconfigD == "" {
			kubeconfigD = os.Getenv("KUBECONFIG_DESTINATION")
		}
		if kubeconfigD == "" {
			fmt.Printf("cannot get destination cluster kube config, please use --destination-kubeconfig or export environment variable KUBECONFIG_DESTINATION to set\n")
			os.Exit(1)
		}

		// For source
		migrationclientS, err := getClient(kubeconfigS, namespaceS)
		if err != nil {
			fmt.Printf(err.Error())
			os.Exit(1)
		}
		err = migrationclientS.PrintServiceWithRevisions("source")
		if err != nil {
			fmt.Printf(err.Error())
			os.Exit(1)
		}

		// For destination
		migrationclientD, err := getClient(kubeconfigD, namespaceD)
		if err != nil {
			fmt.Printf(err.Error())
			os.Exit(1)
		}

		fmt.Println(color.GreenString("[Before migration in destination cluster]"))
		err = migrationclientD.PrintServiceWithRevisions("destination")
		if err != nil {
			fmt.Printf(err.Error())
			os.Exit(1)
		}

		fmt.Println("\nNow migrate all Knative service resources:")
		fmt.Println("From the source", color.BlueString(namespaceS), "namespace of cluster", color.CyanString(kubeconfigS))
		fmt.Println("To the destination", color.BlueString(namespaceD), "namespace of cluster", color.CyanString(kubeconfigD))

		cfgD, err := clientcmd.BuildConfigFromFlags("", kubeconfigS)
		clientset, err := clientset.NewForConfig(cfgD)
		if err != nil {
			fmt.Printf(err.Error())
			os.Exit(1)
		}
		namespaceExists, err := namespaceExists(*clientset, namespaceD)
		if err != nil {
			fmt.Printf(err.Error())
			os.Exit(1)
		}

		if !namespaceExists {
			fmt.Println("Create namespace", color.BlueString(namespaceD), "in destination cluster")
			nsSpec := &apiv1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: namespaceD}}
			_, err = clientset.CoreV1().Namespaces().Create(nsSpec)
			if err != nil {
				fmt.Printf(err.Error())
				os.Exit(1)
			}
		} else {
			fmt.Println("Namespace", namespaceD, "already exists in destination cluster")
		}
		if err != nil {
			fmt.Printf(err.Error())
			os.Exit(1)
		}

		servicesS, err := migrationclientS.ListService()
		if err != nil {
			fmt.Printf(err.Error())
			os.Exit(1)
		}
		for i := 0; i < len(servicesS.Items); i++ {
			serviceS := servicesS.Items[i]
			serviceExists, err := migrationclientD.ServiceExists(serviceS.Name)
			if err != nil {
				fmt.Printf(err.Error())
				os.Exit(1)
			}

			if serviceExists {
				if cmd.Flag("force").Value.String() == "false" {
					fmt.Println("Cannot migrate service", color.CyanString(serviceS.Name), "in namespace", color.BlueString(namespaceS),
						"because the service already exists and no --force option was given\n")
					os.Exit(1)
				}
				fmt.Println("Deleting service", color.CyanString(serviceS.Name), "from the destination cluster and recreate as replacement")
				migrationclientD.DeleteService(serviceS.Name)
				if err != nil {
					fmt.Printf(err.Error())
					os.Exit(1)
				}
			}

			_, err = migrationclientD.CreateService(&serviceS)
			if err != nil {
				fmt.Printf(err.Error())
				os.Exit(1)
			}
			fmt.Println("Migrated service", color.CyanString(serviceS.Name), "Successfully")

			serviceD, err := migrationclientD.GetService(serviceS.Name)
			if err != nil {
				fmt.Printf(err.Error())
				os.Exit(1)
			}

			config, err := migrationclientD.GetConfig(serviceD.Name)
			if err != nil {
				fmt.Printf(err.Error())
				os.Exit(1)
			}
			config_uuid := config.UID

			revisionsS, err := migrationclientS.ListRevisionByService(serviceS.Name)
			if err != nil {
				fmt.Printf(err.Error())
				os.Exit(1)
			}
			for i := 0; i < len(revisionsS.Items); i++ {
				revisionS := revisionsS.Items[i]
				if revisionS.Name != serviceS.Status.LatestReadyRevisionName {
					_, err := migrationclientD.CreateRevision(&revisionS, config_uuid)
					if err != nil {
						fmt.Printf(err.Error())
						os.Exit(1)
					}
					fmt.Println("Migrated revision", color.CyanString(revisionS.Name), "successfully")
				} else {
					retries := 0
					for {
						revision, err := migrationclientD.GetRevision(revisionS.Name)
						if err != nil {
							fmt.Printf(err.Error())
							os.Exit(1)
						}
						sourceRevisionGeneration := revisionS.ObjectMeta.Labels["serving.knative.dev/configurationGeneration"]
						revision.ObjectMeta.Labels["serving.knative.dev/configurationGeneration"] = sourceRevisionGeneration
						err = migrationclientD.UpdateRevision(revision)
						if err != nil {
							// Retry to update when a resource version conflict exists
							if api_errors.IsConflict(err) && retries < MaxUpdateRetries {
								retries++
								continue
							}
							fmt.Printf(err.Error())
							os.Exit(1)
						}
						fmt.Println("Replace revision", color.CyanString(revisionS.Name), "to generation", sourceRevisionGeneration, "successfully")
						break
					}
				}
			}
			fmt.Println("")
		}

		fmt.Println(color.GreenString("[After migration in destination cluster]"))
		err = migrationclientD.PrintServiceWithRevisions("destination")
		if err != nil {
			fmt.Printf(err.Error())
			os.Exit(1)
		}

		if cmd.Flag("delete").Value.String() == "false" {
			fmt.Println("Migrate without --delete option, skip deleting Knative resource in source cluster")
		} else {
			fmt.Println("Migrate with --delete option, deleting all Knative resource in source cluster")
			servicesS, err := migrationclientS.ListService()
			if err != nil {
				fmt.Printf(err.Error())
				os.Exit(1)
			}
			for i := 0; i < len(servicesS.Items); i++ {
				serviceS := servicesS.Items[i]
				err = migrationclientS.DeleteService(serviceS.Name)
				if err != nil {
					fmt.Printf(err.Error())
					os.Exit(1)
				}
				fmt.Println("Deleted service", serviceS.Name, "in source cluster", namespaceS, "namespace")
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)

	migrateCmd.Flags().StringP("namespace", "n", "default", "The namespace of the source Knative resources")

	migrateCmd.Flags().String("destination-kubeconfig", "", "The kubeconfig of the destination Knative resources (default is KUBECONFIG_DESTINATION from environment variable)")
	migrateCmd.Flags().String("destination-namespace", "", "The namespace of the destination Knative resources")

	migrateCmd.Flags().Bool("force", false, "Migrate service forcefully, replaces existing service if any.")
	migrateCmd.Flags().Bool("delete", false, "Delete all Knative resources after kn-migration from source cluster")
}

func getClient(kubeconfig, namespace string) (MigrationClient, error) {
	cfg, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}
	servingClient, err := serving_v1_client.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}
	migrationclient := NewMigrationClient(servingClient, namespace)
	return migrationclient, nil
}

func namespaceExists(client clientset.Clientset, namespace string) (bool, error) {
	_,err := client.CoreV1().Namespaces().Get(namespace, metav1.GetOptions{})
	if api_errors.IsNotFound(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
