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

var (
	// AdminRegistryCmdName is used in the labels to mark the resource this command created
	AdminRegistryCmdName = "kn-admin-registry"
	// DockerJSONName is used to represent ".dockerconfigjson" in secret data
	DockerJSONName = ".dockerconfigjson"
)

// AdminRegistryLabels is a set of labels which will be added to the registry resources to indicate
// that these resources are managed by admin registry command.
var AdminRegistryLabels = map[string]string{
	pkg.LabelManagedBy: AdminRegistryCmdName,
}

// Registry contains data for secret creation
type Registry struct {
	Auths Auths `json:"auths"`
}

// registryCred contains actual credentials which are used to pull images
type registryCred struct {
	Username string `json:"Username"`
	Password string `json:"Password"`
	Email    string `json:"Email"`
}

// Auths is a map of docker credentials indexed by server url
type Auths map[string]registryCred

// NewPrivateRegistryCmd represents the privateRegistry command
func NewPrivateRegistryCmd(p *pkg.AdminParams) *cobra.Command {
	var privateRegistryCmd = &cobra.Command{
		Use:   "registry",
		Short: "Manage registry",
		Long:  `Manage registry used by Knative service deployment`,
	}
	privateRegistryCmd.AddCommand(NewRegistryAddCommand(p))
	privateRegistryCmd.AddCommand(NewRegistryRmCommand(p))
	privateRegistryCmd.InitDefaultHelpCmd()
	return privateRegistryCmd
}
