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
	"github.com/spf13/cobra"
	"knative.dev/client-contrib/plugins/admin/pkg"
)

// NewDomainCmd return the domain root command
func NewDomainCmd(p *pkg.AdminParams) *cobra.Command {
	var domainCmd = &cobra.Command{
		Use:   "domain",
		Short: "Manage route domain",
		Long:  `Manage default route domain or custom route domain for service(s) with selectors`,
	}
	domainCmd.AddCommand(NewDomainSetCommand(p))
	domainCmd.AddCommand(NewDomainUnSetCommand(p))
	domainCmd.InitDefaultHelpCmd()
	return domainCmd
}
