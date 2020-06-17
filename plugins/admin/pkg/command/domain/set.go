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
	"strings"

	"knative.dev/client-contrib/plugins/admin/pkg"

	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var selector []string
var domain string

// NewDomainSetCommand return the command to set knative custom domain
func NewDomainSetCommand(p *pkg.AdminParams) *cobra.Command {
	domainSetCommand := &cobra.Command{
		Use:   "set",
		Short: "set route domain",
		Long: `Set Knative route domain for service

For example:
# To set a default route domain
kn admin domain set --custom-domain mydomain.com
# To set a route domain for service having label app=v1
kn admin domain set --custom-domain mydomain.com --selector app=v1
`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			domain = strings.TrimSpace(domain)
			if domain == "" {
				return errors.New("'domain set' requires the route name provided with the --custom-domain option")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			currentCm := &corev1.ConfigMap{}
			currentCm, err := p.ClientSet.CoreV1().ConfigMaps("knative-serving").Get("config-domain", metav1.GetOptions{})
			if err != nil {
				return fmt.Errorf("failed to get configmaps: %+v", err)
			}
			desiredCm := currentCm.DeepCopy()
			labels := "selector:\n"
			for _, label := range selector {
				k, v, err := splitByEqualSign(label)
				if err != nil {
					return err
				}
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
				if v == value {
					delete(desiredCm.Data, k)
					break
				}
			}

			desiredCm.Data[domain] = value
			if !equality.Semantic.DeepEqual(desiredCm.Data, currentCm.Data) {
				_, err = p.ClientSet.CoreV1().ConfigMaps("knative-serving").Update(desiredCm)
				if err != nil {
					return fmt.Errorf("Failed to update ConfigMaps: %+v", err)
				}
				cmd.Printf("Updated knative route domain to %q\n", domain)
			} else {
				cmd.Printf("Knative route domain %q not changed\n", domain)
			}
			return nil
		},
	}

	domainSetCommand.Flags().StringVarP(&domain, "custom-domain", "d", "", "desired custom domain")
	domainSetCommand.MarkFlagRequired("custom-domain")
	domainSetCommand.Flags().StringSliceVar(&selector, "selector", nil, "domain selector. name=value; you may provide this flag any number of times to set multiple selectors.")
	domainSetCommand.InitDefaultHelpFlag()

	return domainSetCommand
}

func splitByEqualSign(pair string) (string, string, error) {
	parts := strings.Split(pair, "=")
	if len(parts) != 2 || strings.TrimSpace(parts[0]) == "" || strings.TrimSpace(parts[1]) == "" {
		return "", "", fmt.Errorf("expecting the selector format is name=vlue, but given '%s'", pair)
	}
	return strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]), nil
}
