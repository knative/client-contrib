/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"github.com/knative/client-contrib/plugins/kn-image/tekton"
	servingclientset_v1alpha1 "github.com/knative/client/pkg/serving/v1alpha1"
	serving_v1alpha1_api "github.com/knative/serving/pkg/apis/serving/v1alpha1"
	serving_v1beta1_api "github.com/knative/serving/pkg/apis/serving/v1beta1"
	serviceclientset "github.com/knative/serving/pkg/client/clientset/versioned"
	"github.com/spf13/cobra"
	tektoncdclientset "github.com/tektoncd/pipeline/pkg/client/clientset/versioned"
	corev1 "k8s.io/api/core/v1"
	api_errors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc" // from https://github.com/kubernetes/client-go/issues/345
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"time"
)

// deployCmd represents the deploy command
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Build the container image from source code or deploy Knative service",
	Example: `
  # Deploy from git repository to Knative service
  # ( related: https://github.com/knative/client-contrib/blob/master/plugins/kn-image/doc/deploy-git-resource.md )
  kn-imaage deploy cnbtest --builder buildpacks-v3 --git-url https://github.com/zhangtbj/cf-sample-app-nodejs --git-revision master --saved-image us.icr.io/test/cnbtest:v1 --serviceaccount default --force`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("")
		if len(args) < 1 {
			cmd.Help()
			os.Exit(0)
		}
		name := args[0]

		fmt.Println("[INFO] Deploy from git repository to Knative service")
		serviceAccount := cmd.Flag("serviceaccount").Value.String()
		namespace := cmd.Flag("namespace").Value.String()

		// Config kubeconfig
		kubeconfig = cmd.Flag("kubeconfig").Value.String()
		if kubeconfig == "" {
			kubeconfig = os.Getenv("KUBECONFIG")
		}
		cfg, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			fmt.Println("[ERROR] Parsing kubeconfig error:", err)
			os.Exit(1)
		}
		client, err := tektoncdclientset.NewForConfig(cfg)
		if err != nil {
			fmt.Println("[ERROR] Building kubeconfig error:", err)
			os.Exit(1)
		}
		tektonClient := tekton.NewTektonClient(client.TektonV1alpha1(), namespace)

		builder := cmd.Flag("builder").Value.String()
		if builder == "" {
			fmt.Println("[ERROR] Builder cannot be empty, please use --builder to set")
			os.Exit(1)
		}

		gitUrl := cmd.Flag("git-url").Value.String()
		if gitUrl == "" {
			fmt.Println("[ERROR] Git url cannot be empty, please use --git-url to set")
			os.Exit(1)
		}
		gitRevision := cmd.Flag("git-revision").Value.String()
		if gitRevision == "" {
			fmt.Println("[ERROR] Git revision cannot be empty, please use --git-revision to set")
			os.Exit(1)
		}
		image := cmd.Flag("saved-image").Value.String()
		if image == "" {
			fmt.Println("[ERROR] Saved-image cannot be empty, please use --saved-image to set")
			os.Exit(1)
		}

		if len(gitUrl) > 0 {
			err = tektonClient.BuildFromGit(name, builder, gitUrl, gitRevision, image, serviceAccount, namespace)
			if err != nil {
				fmt.Println("[ERROR] Building image error:", err)
				os.Exit(1)
			}
			fmt.Println("[INFO] Generate image", image, "from git repo")
		}

		force, err := cmd.Flags().GetBool("force")
		if err != nil {
			fmt.Println("[ERROR] Parsing force parameter:", err)
			os.Exit(1)
		}
		// Deploy knative service
		knclient, err := serviceclientset.NewForConfig(cfg)
		if err != nil {
			fmt.Println("[ERROR] Serving kubeconfig:", err)
			os.Exit(1)
		}

		fmt.Println("\n[INFO] Deploy the Knative service by using the new generated image")
		servingClient := servingclientset_v1alpha1.NewKnServingClient(knclient.ServingV1alpha1(), namespace)
		serviceExists, service, err := serviceExists(servingClient, name)
		if err != nil {
			fmt.Println("[ERROR] Checking service exist:", err)
			os.Exit(1)
		}
		action := "created"
		if serviceExists {
			if force {
				err = replaceService(servingClient, service, image, serviceAccount)
				action = "replaced"
			} else {
				fmt.Println(
					"[ERROR] cannot create service", name, "in namespace", namespace,
						"because the service already exists and no --force option was given")
				os.Exit(1)
			}
		} else {
			if service == nil {
				service, err = constructService(name, image, serviceAccount, namespace)
			}
			if err != nil {
				fmt.Println("[ERROR] Constructing service:", err)
				os.Exit(1)
			}
			err = createService(servingClient, service)
		}
		if err != nil {
			fmt.Println("[ERROR] Create service:", err)
			os.Exit(1)
		} else {
			fmt.Println("[INFO] Service", service.Name , "successfully", action ,"in namespace", namespace)
		}

		time.Sleep(5 * time.Second)
		i := 0
		for  i < MaxTimeout {
			service, err = servingClient.GetService(name)
			if service.Status.LatestReadyRevisionName != "" {
				fmt.Println("[INFO] service", name,"is ready")
				url := service.Status.URL.String()
				if url == "" {
					url = service.Status.DeprecatedDomain
				}
				fmt.Println("[INFO] Service", name,"url is", url)
				return
			} else {
				fmt.Println("[INFO] Service", name,"is still creating, waiting")
				time.Sleep(5 * time.Second)
			}
			if i == MaxTimeout {
				fmt.Println("[ERROR] Fail to create service", name, "after timeout")
				os.Exit(1)
			}
			i += 5
			time.Sleep(5 * time.Second)
		}
	},
}

func init() {
    rootCmd.AddCommand(deployCmd)

	deployCmd.PersistentFlags().StringP("kubeconfig", "", "", "kube config file (default is KUBECONFIG from ENV property)")
	deployCmd.Flags().StringP("builder", "b", "", "builder of source-to-image task")
	deployCmd.Flags().StringP( "git-url", "u","", "[Git] url of git repo")
	deployCmd.Flags().StringP( "git-revision", "r","master", "[Git] revision of git repo")
	deployCmd.Flags().StringP("saved-image", "i", "", "generated saved image path")
	deployCmd.Flags().StringP("serviceaccount", "s", "default", "service account to push image")
	deployCmd.Flags().StringP( "namespace", "n","default", "namespace of build")
	deployCmd.Flags().BoolP("force", "f", false, "Create service forcefully, replaces existing service if any")
}

// Create a new Knative service
func createService(client servingclientset_v1alpha1.KnClient, service *serving_v1alpha1_api.Service) error {
	err := client.CreateService(service)
	if err != nil {
		return err
	}
	return nil
}

// Replace the existing Knative service
func replaceService(client servingclientset_v1alpha1.KnClient, service *serving_v1alpha1_api.Service, image, serviceAccount string) error {
	var retries = 0
	for {
		service.Spec = serving_v1alpha1_api.ServiceSpec{
			ConfigurationSpec:    serving_v1alpha1_api.ConfigurationSpec{
				Template: &serving_v1alpha1_api.RevisionTemplateSpec{
					Spec: serving_v1alpha1_api.RevisionSpec{
						RevisionSpec: serving_v1beta1_api.RevisionSpec{
							PodSpec: serving_v1beta1_api.PodSpec{
								ServiceAccountName: serviceAccount,
								Containers: []corev1.Container{
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
		err := client.UpdateService(service)
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

// Check if the service exists
func serviceExists(client servingclientset_v1alpha1.KnClient, name string) (bool, *serving_v1alpha1_api.Service, error) {
	service, err := client.GetService(name)
	if api_errors.IsNotFound(err) {
		return false, nil, nil
	}
	if err != nil {
		return false, nil, err
	}
	return true, service, nil
}

// Create service struct from provided options
func constructService(name, image, serviceAccount, namespace string) (*serving_v1alpha1_api.Service,
	error) {

	service := serving_v1alpha1_api.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: serving_v1alpha1_api.ServiceSpec{
			ConfigurationSpec:    serving_v1alpha1_api.ConfigurationSpec{
				Template: &serving_v1alpha1_api.RevisionTemplateSpec{
					Spec: serving_v1alpha1_api.RevisionSpec{
						RevisionSpec: serving_v1beta1_api.RevisionSpec{
							PodSpec: serving_v1beta1_api.PodSpec{
								ServiceAccountName: serviceAccount,
								Containers: []corev1.Container{
									{
										Image: image,
									},
								},
							},
						},
					},
				},
			},
		},
	}

	return &service, nil
}
