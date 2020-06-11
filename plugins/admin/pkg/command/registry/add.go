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
	"knative.dev/client-contrib/plugins/admin/pkg"

	"encoding/json"

	"github.com/spf13/cobra"
)

type prcmdFlags struct {
	Server     string
	SecretName string
	Email      string
	Username   string
	Password   string
}

var prflags prcmdFlags

// NewPrAddCommand represents the add command
func NewPrAddCommand(p *pkg.AdminParams) *cobra.Command {
	var prAddCmd = &cobra.Command{
		Use:   "add",
		Short: "add registry with credentials",
		Long: `add registry with credentials to enable Service deployment from it
For example:

kn admin registry add \
  --secret-name=[SECRET_NAME]
  --server=[REGISTRY_SERVER_URL] \
  --email=[REGISTRY_EMAIL] \
  --username=[REGISTRY_USER] \
  --password=[REGISTRY_PASSWORD]`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if prflags.Username == "" {
				return errors.New("'registry add' requires the registry username to run provided with the --username option")
			}
			if prflags.Password == "" {
				return errors.New("'registry add' requires the registry password to run provided with the --password option")
			}
			if prflags.Server == "" {
				return errors.New("'registry add' requires the registry server to run provided with the --server option")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			dockerCfg := Registry{
				Auths: Auths{
					prflags.Server: DockerCred{
						Username: prflags.Username,
						Password: prflags.Password,
						Email:    prflags.Email,
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
					GenerateName: fmt.Sprintf("%s-", prflags.SecretName),
					Namespace:    "default",
					Labels:       AdminRegistryLabels,
				},
				Data: secretData,
			}

			secret, err = p.ClientSet.CoreV1().Secrets("default").Create(secret)
			if err != nil {
				return fmt.Errorf("failed to create secret: %v", err)
			}

			defaultSa, err := p.ClientSet.CoreV1().ServiceAccounts("default").Get("default", metav1.GetOptions{})
			if err != nil {
				return fmt.Errorf("failed to get serviceaccount: %v", err)
			}
			desiredSa := defaultSa.DeepCopy()
			desiredSa.ImagePullSecrets = append(desiredSa.ImagePullSecrets, corev1.LocalObjectReference{
				Name: secret.Name,
			})

			_, err = p.ClientSet.CoreV1().ServiceAccounts("default").Update(desiredSa)
			if err != nil {
				return fmt.Errorf("failed to add registry secret in default serviceaccount: %v", err)
			}
			cmd.Printf("Private registry %s added for default serviceaccount\n", prflags.Server)
			return nil
		},
	}

	prAddCmd.Flags().StringVar(&prflags.SecretName, "secret-name", "secret-registry", "Registry Secret Name")
	prAddCmd.Flags().StringVar(&prflags.Server, "server", "", "Registry Address")
	prAddCmd.MarkFlagRequired("server")
	prAddCmd.Flags().StringVar(&prflags.Email, "email", "user@default.email.com", "Registry Email")
	prAddCmd.Flags().StringVar(&prflags.Username, "username", "", "Registry Username")
	prAddCmd.MarkFlagRequired("username")
	prAddCmd.Flags().StringVar(&prflags.Password, "password", "", "Registry Password")
	prAddCmd.MarkFlagRequired("password")

	prAddCmd.InitDefaultHelpFlag()
	return prAddCmd
}
