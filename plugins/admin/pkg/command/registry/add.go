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
	"fmt"
	"os"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"knative.dev/client-contrib/plugins/admin/pkg"

	"encoding/json"

	"github.com/spf13/cobra"
)

type prcmdFlags struct {
	DockerServer   string
	SecretName     string
	DockerEmail    string
	DockerUsername string
	DockerPassword string
}

type DockerRegistry struct {
	Auths Auths `json:"auths"`
}
type RegistryCred struct {
	Username string `json:"Username"`
	Password string `json:"Password"`
	Email    string `json:"Email"`
}

//type Auths struct {
//	RegistryCred RegistryCred `json:"us.icr.io"`
//}
type Auths map[string]RegistryCred

//
//type Info map[string]Person

var prflags prcmdFlags

// addCmd represents the add command

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
		Run: func(cmd *cobra.Command, args []string) {
			dockerCfg := DockerRegistry{
				Auths: Auths{
					prflags.DockerServer: RegistryCred{
						Username: prflags.DockerUsername,
						Password: prflags.DockerPassword,
						Email:    prflags.DockerEmail,
					},
				},
			}

			j, err := json.Marshal(dockerCfg)
			if err != nil {
				panic(err)
			}

			secretData := map[string][]byte{
				".dockerconfigjson": j,
			}

			secret := &corev1.Secret{
				TypeMeta: metav1.TypeMeta{
					Kind:       "Secret",
					APIVersion: "v1",
				},
				Type: corev1.SecretTypeDockerConfigJson,
				ObjectMeta: metav1.ObjectMeta{
					Name:      prflags.SecretName,
					Namespace: "default",
				},
				Data: secretData,
			}

			_, err = p.ClientSet.CoreV1().Secrets("default").Create(secret)
			if err != nil {
				fmt.Println("failed to create secret:", err)
				os.Exit(1)
			}

			defaultSa, err := p.ClientSet.CoreV1().ServiceAccounts("default").Get("default", metav1.GetOptions{})
			desiredSa := defaultSa.DeepCopy()
			desiredSa.ImagePullSecrets = []corev1.LocalObjectReference{{
				Name: prflags.SecretName,
			}}

			_, err = p.ClientSet.CoreV1().ServiceAccounts("default").Update(desiredSa)
			if err != nil {
				fmt.Println("Failed to add registry secret in default Service Account:\n", err)
				os.Exit(1)
			}

			fmt.Printf("Private registry %s added for default Service Account\n", prflags.DockerServer)
		},
	}

	prAddCmd.Flags().StringVar(&prflags.SecretName, "secret-name", "", "Registry Secret Name")
	prAddCmd.Flags().StringVar(&prflags.DockerServer, "server", "", "Registry Address")
	prAddCmd.Flags().StringVar(&prflags.DockerEmail, "email", "", "Registry Email")
	prAddCmd.Flags().StringVar(&prflags.DockerUsername, "username", "", "Registry Username")
	prAddCmd.Flags().StringVar(&prflags.DockerPassword, "password", "", "Registry Email")

	return prAddCmd
}
