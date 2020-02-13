/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

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
package domain

import (
	"fmt"
	"github.com/knative/client-contrib/plugins/kn-admin/cmd"
	"github.com/spf13/cobra"
)

// domainCmd represents the domain command
func NewDomainCmd(p *cmd.AdminParams) *cobra.Command {
	var domainCmd = &cobra.Command{
		Use:   "domain",
		Short: "Manage route domain",
		Long: `List and set default route domain or route domain for Service with selectors. For example:

kn admin domain list - to list Knative route domain
kn admin domain set - to set Knative route domain`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("domain called")
		},
	}
	domainCmd.AddCommand(NewDomainSetCommand(p))
	return domainCmd
}
