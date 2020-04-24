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

package list

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc" // from https://github.com/kubernetes/client-go/issues/345
	"k8s.io/client-go/tools/clientcmd"
	"knative.dev/client-contrib/plugins/migration/pkg/command"
	servingv1client "knative.dev/serving/pkg/client/clientset/versioned/typed/serving/v1"
)

type listCmdFlags struct {
	Namespace  string
	KubeConfig string
}

var listFlags listCmdFlags

// listCmd represents the list command
func NewListCommand() *cobra.Command {
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all Knative service resources",
		Long:  `List all Knative service resources`,
		Run: func(cmd *cobra.Command, args []string) {
			kubeConfig := listFlags.KubeConfig
			if kubeConfig == "" {
				kubeConfig = os.Getenv("KUBECONFIG")
			}
			ServingClient, err := getClient(kubeConfig, listFlags.Namespace)
			if err != nil {
				fmt.Errorf(err.Error())
				os.Exit(1)
			}
			err = ServingClient.PrintServiceWithRevisions("current")
			if err != nil {
				fmt.Errorf(err.Error())
				os.Exit(1)
			}
		},
	}

	listCmd.Flags().StringVarP(&listFlags.Namespace, "namespace", "n", "default", "The namespace of the Knative resources (default is default namespace)")
	listCmd.Flags().StringVar(&listFlags.KubeConfig, "kubeconfig", "", "The kubeconfig of the Knative resources (default is KUBECONFIG from environment variable)")
	return listCmd
}

func getClient(kubeConfig, namespace string) (command.MigrationClient, error) {
	cfg, err := clientcmd.BuildConfigFromFlags("", kubeConfig)
	if err != nil {
		return nil, err
	}
	servingClient, err := servingv1client.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}
	migrationClient := command.NewMigrationClient(servingClient, namespace)
	return migrationClient, nil
}
