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
	"github.com/knative/client-contrib/plugins/kn-image/clients"
	serviceclientset "github.com/knative/serving/pkg/client/clientset/versioned"
	"github.com/spf13/cobra"
	tektoncdclientset "github.com/tektoncd/pipeline/pkg/client/clientset/versioned"
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
		tektonClient := clients.NewTektonClient(client.TektonV1alpha1(), namespace)

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
		servingClient := clients.NewServingClient(knclient.ServingV1(), namespace)
		serviceExists, service, err := servingClient.ServiceExists(name)
		if err != nil {
			fmt.Println("[ERROR] Checking service exist:", err)
			os.Exit(1)
		}
		action := "created"
		if serviceExists {
			if force {
				err = servingClient.UpdateService(service, image, serviceAccount)
				action = "replaced"
			} else {
				fmt.Println(
					"[ERROR] cannot create service", name, "in namespace", namespace,
						"because the service already exists and no --force option was given")
				os.Exit(1)
			}
		} else {
			if service == nil {
				service = servingClient.ConstructService(name, image, serviceAccount, namespace)
			}
			if err != nil {
				fmt.Println("[ERROR] Constructing service:", err)
				os.Exit(1)
			}
			err = servingClient.CreateService(service)
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
					url = service.Status.URL.String()
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
