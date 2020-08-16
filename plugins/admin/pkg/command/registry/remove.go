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
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"

	"knative.dev/client-contrib/plugins/admin/pkg"
	"knative.dev/client/pkg/kn/commands"

	"github.com/spf13/cobra"
)

var username string
var server string
var serviceaccount string

// NewRegistryRmCommand represents the remove command
func NewRegistryRmCommand(p *pkg.AdminParams) *cobra.Command {
	var registryRmCmd = &cobra.Command{
		Use:     "remove",
		Aliases: []string{"rm"},
		Short:   "Remove registry settings",
		Long:    `Remove registry settings by server and username to delete secret and update ServiceAccount`,
		Example: `
  # To remove registry settings
  kn admin registry remove \
    --username=[REGISTRY_USER] \
    --server=[REGISTRY_SERVER_URL] \
    --namespace=[NAMESPACE] \
    --serviceaccount=[SERVICE_ACCOUNT]`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if username == "" {
				return errors.New("'registry remove' requires the registry username provided with the --username option")
			}
			if server == "" {
				return errors.New("'registry remove' requires the registry server url provided with the --server option")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			namespace := cmd.Flag("namespace").Value.String()
			if namespace == "" {
				namespace = "default"
			}
			// get all credential secrets which have the label managed-by=kn-admin-registry
			secrets, err := p.ClientSet.CoreV1().Secrets(namespace).List(metav1.ListOptions{
				LabelSelector: labels.SelectorFromSet(AdminRegistryLabels).String(),
			})
			if err != nil {
				return fmt.Errorf("failed to list secret: %v", err)
			}

			// filter the secrets with username and server
			secretsMap := make(map[string]corev1.Secret)
			for _, secret := range secrets.Items {
				registry := Registry{}
				err = json.Unmarshal(secret.Data[DockerJSONName], &registry)
				if err != nil {
					return fmt.Errorf("failed unmarshal secret data '.dockerconfigjson': %v", err)
				}
				for secretServer, secretAuth := range registry.Auths {
					if secretServer == server && secretAuth.Username == username {
						secretsMap[secret.Name] = *secret.DeepCopy()
					}
				}
			}
			if len(secretsMap) == 0 {
				cmd.Printf("No registry found for server: '%s' and username: '%s'\n", server, username)
				return nil
			}

			sa, err := p.ClientSet.CoreV1().ServiceAccounts(namespace).Get(serviceaccount, metav1.GetOptions{})
			if err != nil {
				return fmt.Errorf("failed to get serviceaccount '%s' in namespace '%s': %v", serviceaccount, namespace, err)
			}

			desiredSa := sa.DeepCopy()
			imagePullSecrets := []corev1.LocalObjectReference{}
			for _, ips := range desiredSa.ImagePullSecrets {
				if _, ok := secretsMap[ips.Name]; !ok {
					// only store the secrets that do not exist in the map
					imagePullSecrets = append(imagePullSecrets, ips)
				}

			}

			desiredSa.ImagePullSecrets = imagePullSecrets
			_, err = p.ClientSet.CoreV1().ServiceAccounts(namespace).Update(desiredSa)
			if err != nil {
				return fmt.Errorf("failed to remove registry secret in serviceaccount '%s' in namespace '%s': %v", serviceaccount, namespace, err)
			}
			cmd.Printf("ImagePullSecrets of serviceaccount '%s' in namespace '%s' is updated\n", desiredSa.Name, desiredSa.Namespace)

			deleteSecretsErrCh := make(chan error, len(secretsMap))
			deleteSecrets(cmd, p.ClientSet, secretsMap, deleteSecretsErrCh)

			var deleteSecretsErr error
			select {
			case deleteSecretsErr = <-deleteSecretsErrCh:
				if deleteSecretsErr != nil {
					break
				}
			default:
			}

			if deleteSecretsErr != nil {
				return fmt.Errorf("failed to delete secrets: %v", deleteSecretsErr)
			}

			return nil
		},
	}

	commands.AddNamespaceFlags(registryRmCmd.Flags(), false)
	registryRmCmd.Flags().StringVar(&serviceaccount, "serviceaccount", "default", "the service account to save imagePullSecrets")
	registryRmCmd.Flags().StringVar(&username, "username", "", "registry username")
	registryRmCmd.MarkFlagRequired("username")
	registryRmCmd.Flags().StringVar(&server, "server", "", "registry address")
	registryRmCmd.MarkFlagRequired("server")
	registryRmCmd.InitDefaultHelpFlag()
	return registryRmCmd
}

func deleteSecrets(cmd *cobra.Command, clientset kubernetes.Interface, secretsMap map[string]corev1.Secret, errCh chan<- error) {
	w := sync.WaitGroup{}
	w.Add(len(secretsMap))
	for _, s := range secretsMap {
		go func(secret corev1.Secret) {
			defer w.Done()
			err := clientset.CoreV1().Secrets(secret.Namespace).Delete(secret.Name, &metav1.DeleteOptions{})
			if err != nil {
				if apierrors.IsNotFound(err) {
					cmd.Printf("Secret '%s' in namespace '%s' is not found, skipped\n", secret.Name, secret.Namespace)
				} else {
					errCh <- fmt.Errorf("failed to delete secret '%s' in namespace '%s': %v", secret.Name, secret.Namespace, err)
				}
			} else {
				cmd.Printf("Secret '%s' in namespace '%s' is deleted\n", secret.Name, secret.Namespace)
			}
		}(s)
	}
	w.Wait()
}
