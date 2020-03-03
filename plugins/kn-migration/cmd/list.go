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
	"github.com/spf13/cobra"
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc" // from https://github.com/kubernetes/client-go/issues/345
	"k8s.io/client-go/tools/clientcmd"
	serving_v1_client "knative.dev/serving/pkg/client/clientset/versioned/typed/serving/v1"
	"os"
)

var kubeconfig string

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all knative service resources",
	Long: `List all knative service resources`,
	Run: func(cmd *cobra.Command, args []string) {
		kubeconfig = rootCmd.Flag("kubeconfig").Value.String()
		if kubeconfig == "" {
			kubeconfig = os.Getenv("KUBECONFIG")
		}
		cfg, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			fmt.Printf(err.Error())
			os.Exit(1)
		}
		serving_client, err := serving_v1_client.NewForConfig(cfg)
		if err != nil {
			fmt.Printf(err.Error())
			os.Exit(1)
		}
		namespace := cmd.Flag("namespace").Value.String()
		migrationClient := NewMigrationClient(serving_client, namespace)
		err = migrationClient.PrintServiceWithRevisions("current")
		if err != nil {
			fmt.Errorf(err.Error())
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().StringP("namespace", "n","default", "The namespace of the knative resources (default is KUBECONFIG from ENV property)")
}
