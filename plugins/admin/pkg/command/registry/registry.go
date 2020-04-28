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
	"github.com/spf13/cobra"
	"knative.dev/client-contrib/plugins/admin/pkg"
)

// privateRegistryCmd represents the privateRegistry command
func NewPrivateRegistryCmd(p *pkg.AdminParams) *cobra.Command {
	var privateRegistryCmd = &cobra.Command{
		Use:   "registry",
		Short: "Manage registry",
		Long: `Manage Service deployment from a registry with credentials
For example:

kn admin registry add \
  --secret-name=[SECRET_NAME]
  --server=[REGISTRY_SERVER_URL] \
  --email=[REGISTRY_EMAIL] \
  --username=[REGISTRY_USER] \
  --password=[REGISTRY_PASSWORD]`,
	}
	privateRegistryCmd.AddCommand(NewPrAddCommand(p))
	return privateRegistryCmd
}
