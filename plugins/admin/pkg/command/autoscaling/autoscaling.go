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

package autoscaling

import (
	"github.com/spf13/cobra"
	"knative.dev/client-contrib/plugins/admin/pkg"
)

// domainCmd represents the domain command
func NewAutoscalingCmd(p *pkg.AdminParams) *cobra.Command {
	var AutoscalingCmd = &cobra.Command{
		Use:   "autoscaling",
		Short: "Manage autoscaling config",
		Long:  `Manage autoscaling provided by Knative Pod Autoscaler (KPA)`,
	}
	AutoscalingCmd.AddCommand(NewAutoscalingUpdateCommand(p))
	AutoscalingCmd.InitDefaultHelpFlag()
	return AutoscalingCmd
}
