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
	"errors"
	"fmt"

	"knative.dev/client-contrib/plugins/admin/pkg/command/utils"

	"knative.dev/client-contrib/plugins/admin/pkg"

	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"

	"knative.dev/client/pkg/kn/flags"
)

var (
	scaleToZero       bool
	enableScaleToZero = "enable-scale-to-zero"
	knativeServing    = "knative-serving"
	configAutoscaler  = "config-autoscaler"
)

func NewAutoscalingUpdateCommand(p *pkg.AdminParams) *cobra.Command {
	AutoscalingUpdateCommand := &cobra.Command{
		Use:   "update",
		Short: "Update autoscaling config",
		Long:  `Update autoscaling config provided by Knative Pod Autoscaler (KPA)`,
		Example: `
  # To enable scale-to-zero
  kn admin autoscaling update --scale-to-zero

  # To disable scale-to-zero
  kn admin autoscaling update --no-scale-to-zero`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if cmd.Flags().NFlag() == 0 {
				return errors.New("'autoscaling update' requires flag(s)")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			var desiredScaleToZero string
			if cmd.Flags().Changed("scale-to-zero") {
				desiredScaleToZero = "true"
			} else if cmd.Flags().Changed("no-scale-to-zero") {
				desiredScaleToZero = "false"
			}

			currentCm := &corev1.ConfigMap{}
			currentCm, err := p.ClientSet.CoreV1().ConfigMaps(knativeServing).Get(configAutoscaler, metav1.GetOptions{})
			if err != nil {
				return fmt.Errorf("failed to get ConfigMaps: %+v", err)
			}
			desiredCm := currentCm.DeepCopy()
			desiredCm.Data[enableScaleToZero] = desiredScaleToZero

			err = utils.UpdateConfigMap(p.ClientSet, desiredCm)
			if err != nil {
				return fmt.Errorf("failed to update ConfigMap %s in namespace %s: %+v", configAutoscaler, knativeServing, err)
			}
			cmd.Printf("Updated Knative autoscaling config %s: %s\n", enableScaleToZero, desiredScaleToZero)

			return nil
		},
	}

	flags.AddBothBoolFlagsUnhidden(AutoscalingUpdateCommand.Flags(), &scaleToZero, "scale-to-zero", "", true,
		"Enable scale-to-zero if set.")

	return AutoscalingUpdateCommand
}
