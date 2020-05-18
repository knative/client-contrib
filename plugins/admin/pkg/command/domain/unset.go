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

package domain

import (
	"errors"
	"fmt"
	"os"

	"knative.dev/client-contrib/plugins/admin/pkg"

	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
)

// updateCmd represents the update command

func NewDomainUnSetCommand(p *pkg.AdminParams) *cobra.Command {
	domainUnSetCommand := &cobra.Command{
		Use:   "unset",
		Short: "unset route domain",
		Long: `unset Knative route domain for service

For example:
# To unset a route domain
kn admin domain unset --custom-domain mydomain.com
`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if domain == "" {
				return errors.New("'domain unset' requires the route name to run provided with the --custom-domain option")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			currentCm := &corev1.ConfigMap{}
			currentCm, err := p.ClientSet.CoreV1().ConfigMaps("knative-serving").Get("config-domain", metav1.GetOptions{})
			if err != nil {
				fmt.Println("failed to get ConfigMaps:", err)
				os.Exit(1)
			}

			desiredCm := currentCm.DeepCopy()

			_, ok := desiredCm.Data[domain]
			if ok {
				delete(desiredCm.Data, domain)
			} else {
				fmt.Printf("Knative route domain %s not found\n", domain)
				os.Exit(1)
			}

			_, err = p.ClientSet.CoreV1().ConfigMaps("knative-serving").Update(desiredCm)
			if err != nil {
				fmt.Println("failed to update ConfigMaps:", err)
				os.Exit(1)
			}
			fmt.Printf("Unset Knative route domain %s\n", domain)
		},
	}

	domainUnSetCommand.Flags().StringVarP(&domain, "custom-domain", "d", "", "custom domain to unset")
	domainUnSetCommand.MarkFlagRequired("custom-domain")

	domainUnSetCommand.InitDefaultHelpFlag()

	return domainUnSetCommand
}
