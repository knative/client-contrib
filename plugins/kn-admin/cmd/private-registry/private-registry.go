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
	"github.com/knative/client-contrib/plugins/kn-admin/cmd"
	"github.com/spf13/cobra"
)

// privateRegistryCmd represents the privateRegistry command
func NewPrivateRegistryCmd(p *cmd.AdminParams) *cobra.Command {
	var privateRegistryCmd = &cobra.Command{
		Use:   "private-registry",
		Short: "Manage private-registry",
		Long: `Manage Service deployment from a private registry
For example:

kn admin private-registry enable \
  --secret-name=[SECRET_NAME]
  --docker-server=[PRIVATE_REGISTRY_SERVER_URL] \
  --docker-email=[PRIVATE_REGISTRY_EMAIL] \
  --docker-username=[PRIVATE_REGISTRY_USER] \
  --docker-password=[PRIVATE_REGISTRY_PASSWORD]`,
	}
	privateRegistryCmd.AddCommand(NewPrEnableCommand(p))
	return privateRegistryCmd
}
