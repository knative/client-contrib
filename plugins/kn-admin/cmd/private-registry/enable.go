/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

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
package private_registry

import (
	"fmt"
	"github.com/knative/client-contrib/plugins/kn-admin/cmd"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"

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
type UsIcrIo struct {
	Username string `json:"Username"`
	Password string `json:"Password"`
	Email    string `json:"Email"`
}
type Auths struct {
	UsIcrIo UsIcrIo `json:"us.icr.io"`
}

var prflags prcmdFlags

// enableCmd represents the enable command

func NewPrEnableCommand(p *cmd.AdminParams) *cobra.Command {
	var prEnableCmd = &cobra.Command{
		Use:   "enable",
		Short: "enable Service deployment from a private registry",
		Long: `enable Service deployment from a private registry
For example:

kn admin private-registry enable \
  --secret-name=[SECRET_NAME]
  --docker-server=[PRIVATE_REGISTRY_SERVER_URL] \
  --docker-email=[PRIVATE_REGISTRY_EMAIL] \
  --docker-username=[PRIVATE_REGISTRY_USER] \
  --docker-password=[PRIVATE_REGISTRY_PASSWORD]`,
		Run: func(cmd *cobra.Command, args []string) {
			dockerCfg := DockerRegistry{
				Auths: Auths{
					UsIcrIo: UsIcrIo{
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

			fmt.Printf("Private registry %s enabled for default Service Account\n", prflags.DockerServer)
		},
	}

	prEnableCmd.Flags().StringVar(&prflags.SecretName, "secret-name", "", "Registry Secret Name")
	prEnableCmd.Flags().StringVar(&prflags.DockerServer, "docker-server", "", "Registry Address")
	prEnableCmd.Flags().StringVar(&prflags.DockerEmail, "docker-email", "", "Registry Email")
	prEnableCmd.Flags().StringVar(&prflags.DockerUsername, "docker-username", "", "Registry Username")
	prEnableCmd.Flags().StringVar(&prflags.DockerPassword, "docker-password", "", "Registry Email")

	return prEnableCmd
}
