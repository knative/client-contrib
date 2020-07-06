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

	"knative.dev/client-contrib/plugins/admin/pkg"
	"knative.dev/client-contrib/plugins/admin/pkg/command/utils"

	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
)

func NewDomainUnSetCommand(p *pkg.AdminParams) *cobra.Command {
	domainUnSetCommand := &cobra.Command{
		Use:   "unset",
		Short: "Unset route domain",
		Long:  `Unset Knative route domain for service(s)`,
		Example: `
  # To unset a route domain
  kn admin domain unset --custom-domain mydomain.com`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if domain == "" {
				return errors.New("'domain unset' requires the route name to run provided with the --custom-domain option")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			currentCm := &corev1.ConfigMap{}
			currentCm, err := p.ClientSet.CoreV1().ConfigMaps(knativeServing).Get(configDomain, metav1.GetOptions{})
			if err != nil {
				return fmt.Errorf("failed to get configmaps: %+v", err)
			}

			desiredCm := currentCm.DeepCopy()

			_, ok := desiredCm.Data[domain]
			if ok {
				delete(desiredCm.Data, domain)
			} else {
				return fmt.Errorf("Knative route domain %s not found\n", domain)
			}

			err = utils.UpdateConfigMap(p.ClientSet, desiredCm)
			if err != nil {
				return fmt.Errorf("failed to update ConfigMap %s in namespace %s: %+v", configDomain, knativeServing, err)
			}

			cmd.Printf("Unset Knative route domain %s\n", domain)
			return nil
		},
	}

	domainUnSetCommand.Flags().StringVarP(&domain, "custom-domain", "d", "", "custom domain to unset")
	domainUnSetCommand.MarkFlagRequired("custom-domain")

	domainUnSetCommand.InitDefaultHelpFlag()

	return domainUnSetCommand
}
