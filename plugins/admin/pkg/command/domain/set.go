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
	"strings"

	"knative.dev/client-contrib/plugins/admin/pkg"

	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
)

var selector []string
var domain string

// updateCmd represents the update command

func NewDomainSetCommand(p *pkg.AdminParams) *cobra.Command {
	domainSetCommand := &cobra.Command{
		Use:   "set",
		Short: "set route domain",
		Long: `set Knative route domain for service

For example:
# To set a default route domain
kn admin domain set --custom-domain mydomain.com
# To set a route domain for service having label app=v1
kn admin domain set --custom-domain mydomain.com --selector app=v1
`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if domain == "" {
				return errors.New("'domain set' requires the route name to run provided with the --custom-domain option")
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
			labels := "selector:\n"
			for _, label := range selector {
				k, v, _ := splitByEqualSign(label)
				label = fmt.Sprintf("  %s: %s\n", k, v)
				labels += label
			}

			var value string
			if len(selector) == 0 {
				value = ""
			} else {
				value = labels
			}

			for k, v := range desiredCm.Data {
				if k == domain && v == "" {
					delete(desiredCm.Data, k)
					break
				}
				if v == value {
					delete(desiredCm.Data, k)
					break
				}
			}

			desiredCm.Data[domain] = value
			if !equality.Semantic.DeepEqual(desiredCm.Data, currentCm.Data) {
				_, err = p.ClientSet.CoreV1().ConfigMaps("knative-serving").Update(desiredCm)
				if err != nil {
					fmt.Println("failed to update ConfigMaps:", err)
					os.Exit(1)
				}
				fmt.Printf("Updated Knative route domain %s\n", domain)
			} else {
				fmt.Printf("Knative route domain %s not changed\n", domain)
			}
		},
	}

	domainSetCommand.Flags().StringVarP(&domain, "custom-domain", "d", "", "Desired custom domain")
	domainSetCommand.MarkFlagRequired("custom-domain")
	domainSetCommand.Flags().StringSliceVar(&selector, "selector", nil, "Domain selector")

	domainSetCommand.InitDefaultHelpFlag()

	return domainSetCommand
}

func splitByEqualSign(pair string) (string, string, error) {
	parts := strings.Split(pair, "=")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("expecting the value format in value1=value2, given %s", pair)
	}
	return parts[0], strings.TrimSuffix(parts[1], "%"), nil
}
