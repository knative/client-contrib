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

package registry

import (
	"errors"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"encoding/json"

	"github.com/spf13/cobra"
)

type registrycmdFlags struct {
	Server     string
	SecretName string
	Email      string
	Username   string
	Password   string
}

var registryFlags registrycmdFlags

// NewRegistryAddCommand represents the add command
func NewRegistryAddCommand(p *registryAdminParams) *cobra.Command {
	var registryAddCmd = &cobra.Command{
		Use:   "add",
		Short: "Add registry with credentials",
		Long:  `Add registry with credentials and enable service deployments from this registry`,
		Example: `
  # add registry with credentials
  kn admin registry add \
    --secret=[SECRET_NAME] \
    --server=[REGISTRY_SERVER_URL] \
    --email=[REGISTRY_EMAIL] \
    --username=[REGISTRY_USER] \
    --password=[REGISTRY_PASSWORD] \
    --namespace=[NAMESPACE] \
    --serviceaccount=[SERVICE_ACCOUNT]`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if registryFlags.Username == "" {
				return errors.New("'registry add' requires the registry username to run provided with the --username option")
			}
			if registryFlags.Password == "" {
				return errors.New("'registry add' requires the registry password to run provided with the --password option")
			}
			if registryFlags.Server == "" {
				return errors.New("'registry add' requires the registry server to run provided with the --server option")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if p.Namespace == "" {
				cmd.Print("No namespace specified, using 'default' namespace\n")
				p.Namespace = "default"
			}
			if p.ServiceAccount == "" {
				cmd.Print("No serviceaccount specified, using 'default' serviceaccount\n")
				p.ServiceAccount = "default"
			}
			dockerCfg := Registry{
				Auths: Auths{
					registryFlags.Server: registryCred{
						Username: registryFlags.Username,
						Password: registryFlags.Password,
						Email:    registryFlags.Email,
					},
				},
			}

			j, err := json.Marshal(dockerCfg)

			secretData := map[string][]byte{
				DockerJSONName: j,
			}

			secret := &corev1.Secret{
				TypeMeta: metav1.TypeMeta{
					Kind:       "Secret",
					APIVersion: "v1",
				},
				Type: corev1.SecretTypeDockerConfigJson,
				ObjectMeta: metav1.ObjectMeta{
					GenerateName: fmt.Sprintf("%s-", registryFlags.SecretName),
					Namespace:    p.Namespace,
					Labels:       AdminRegistryLabels,
				},
				Data: secretData,
			}

			secret, err = p.AdminParams.ClientSet.CoreV1().Secrets(p.Namespace).Create(secret)
			if err != nil {
				return fmt.Errorf("failed to create secret in namespace '%s': %v", p.Namespace, err)
			}

			sa, err := p.AdminParams.ClientSet.CoreV1().ServiceAccounts(p.Namespace).Get(p.ServiceAccount, metav1.GetOptions{})
			if err != nil {
				return fmt.Errorf("failed to get serviceaccount '%s' in namespace '%s': %v", p.ServiceAccount, p.Namespace, err)
			}
			desiredSa := sa.DeepCopy()
			desiredSa.ImagePullSecrets = append(desiredSa.ImagePullSecrets, corev1.LocalObjectReference{
				Name: secret.Name,
			})

			_, err = p.AdminParams.ClientSet.CoreV1().ServiceAccounts(p.Namespace).Update(desiredSa)
			if err != nil {
				return fmt.Errorf("failed to add registry secret in serviceaccount '%s' in namespace '%s': %v", p.ServiceAccount, p.Namespace, err)
			}
			cmd.Printf("Private registry '%s' is added for serviceaccount '%s' in namespace '%s'\n", registryFlags.Server, p.ServiceAccount, p.Namespace)
			return nil
		},
	}

	registryAddCmd.Flags().StringVar(&registryFlags.SecretName, "secret", "secret-registry", "registry secret name")
	registryAddCmd.Flags().StringVar(&registryFlags.Server, "server", "", "registry address")
	registryAddCmd.MarkFlagRequired("server")
	registryAddCmd.Flags().StringVar(&registryFlags.Email, "email", "user@default.email.com", "registry email")
	registryAddCmd.Flags().StringVar(&registryFlags.Username, "username", "", "registry username")
	registryAddCmd.MarkFlagRequired("username")
	registryAddCmd.Flags().StringVar(&registryFlags.Password, "password", "", "registry password")
	registryAddCmd.MarkFlagRequired("password")

	registryAddCmd.InitDefaultHelpFlag()
	return registryAddCmd
}
