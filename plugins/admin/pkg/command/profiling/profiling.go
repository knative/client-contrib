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

package profiling

import (
	//"encoding/json"
	//"errors"
	"fmt"
	"strconv"
	"strings"
	//"sync"

	"github.com/spf13/cobra"

	corev1 "k8s.io/api/core/v1"
	//apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	//"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"knative.dev/client-contrib/plugins/admin/pkg"
	"knative.dev/client-contrib/plugins/admin/pkg/command/utils"
)

const (
	profilingExample = `
  # To enable profiling 
  kn admin profiling --enable

  # To download heap profile data of autoscaler
  kn admin profiling --target autoscaler --heap

  # To download 2 minutes execution trace data of networking-istio
  kn admin profiling --target networking-istio --trace 2m

  # To download cpu profile data of activator and save to specified folder
  kn admin profiling --target activator --profile --save-path /tmp

  # To download all available profile data for specified pod activator-5979f56548
  kn admin profiling --target activator-5979f56548 --all
`

	targetFlagUsgae       = "The profiling target. It can be a Knative component name or a specific pod name, e.g: activator or activator-5979f56548"
	savePathFlagUsage     = "The path to save the downloaded profile data, if not speicifed, the data will be saved in current working folder"
	heapFlagUsage         = "Download heap profile data"
	cpuFlagUsage          = "Download cpu profile data, you can specify the duration with seconds, minutes or hours, e.g: 1m for one minute cpu profile data"
	blockFlagUsgae        = "Download Go routine blocking data"
	traceFlagUsage        = "Download execution trace data, you can specify the duration with seconds, minutes or hours, e.g: 5s for 5 seconds trace data"
	memoryAllocsFlagUsage = "Download all memory allocations data"
	mutexFlagUsage        = "Download holders of contended mutexes data"
	goroutineFlagUsage    = "Download stack traces of all current goroutines"
	threadCreateFlagUsage = "Download stack traces that led to the creation of new OS threads"

	knNamespace  = "knative-serving"
	obsConfigMap = "config-observability"
)

// NewProfilingCommand creates a profiling command
func NewProfilingCommand(p *pkg.AdminParams) *cobra.Command {
	var enableProfiling bool
	var disableProfiling bool
	var target string
	var savePath string
	//var heapFlagVal bool
	var cpuFlagVal string
	//var blockFlagVal bool
	var traceFlagVal string
	//var memAllocsFlagVal bool
	//var goroutineFlagVal bool
	//var threadCreateFlagVal bool

	var profilingCmd = &cobra.Command{
		Use:     "profiling",
		Aliases: []string{"prof"},
		Short:   "Profiling Knative components",
		Long:    `Enable Knative components profiling and download profile data`,
		Example: profilingExample,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if cmd.Flags().Changed("enable") {
				return configProfiling(p.ClientSet, cmd, true)
			} else if cmd.Flags().Changed("disable") {
				return configProfiling(p.ClientSet, cmd, false)
			}
			return nil
		},
	}

	flags := profilingCmd.Flags()
	flags.BoolVar(&enableProfiling, "enable", false, "Enable Knative profiling")
	flags.BoolVar(&disableProfiling, "disable", false, "Disable Knative profiling")
	flags.StringVarP(&target, "target", "t", "", targetFlagUsgae)
	flags.StringVarP(&savePath, "save-path", "s", "", savePathFlagUsage)
	flags.Bool("heap", false, heapFlagUsage)
	flags.StringVar(&cpuFlagVal, "cpu", "30s", cpuFlagUsage)
	flags.Bool("block", false, blockFlagUsgae)
	flags.StringVar(&traceFlagVal, "trace", "30s", traceFlagUsage)
	flags.Bool("mem-allocs", false, memoryAllocsFlagUsage)
	flags.Bool("mutex", false, mutexFlagUsage)
	flags.Bool("goroutine", false, goroutineFlagUsage)
	flags.Bool("thread-create", false, threadCreateFlagUsage)
	return profilingCmd
}

func configProfiling(c kubernetes.Interface, cmd *cobra.Command, enable bool) error {
	currentCm := &corev1.ConfigMap{}
	currentCm, err := c.CoreV1().ConfigMaps(knNamespace).Get(obsConfigMap, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to get ConfigMap %s in namespace %s: %+v", obsConfigMap, knNamespace, err)
	}

	desiredCm := currentCm.DeepCopy()
	if enable {
		desiredCm.Data["profiling.enable"] = "true"
	} else {
		desiredCm.Data["profiling.enable"] = "false"
	}

	//cmd.Printf("%+v", *desiredCm)
	//cmd.Printf("Enable: %+v", currentCm.Data["profiling.enable"])
	err = utils.UpdateConfigMap(c, desiredCm)
	if err != nil {
		return fmt.Errorf("failed to update ConfigMap %s in namespace %s: %+v", obsConfigMap, knNamespace, err)
	}

	if enable {
		cmd.Print("Knative profiling is enabled")
	} else {
		cmd.Print("Knative profiling is disabled")
	}
	return nil
}

func isProfilingEnabled(c kubernetes.Interface) (bool, error) {
	currentCm := &corev1.ConfigMap{}
	currentCm, err := c.CoreV1().ConfigMaps(knNamespace).Get(obsConfigMap, metav1.GetOptions{})
	if err != nil {
		return false, fmt.Errorf("failed to get ConfigMap %s in namespace %s: %+v", obsConfigMap, knNamespace, err)
	}

	if strings.ToLower(currentCm.Data["profiling.enable"]) == "true" {
		return true, nil
	}
	return false, nil
}

// parseDuration parses the given duration string to integer seconds, the duration is an integer plusing a character
// 's', 'm' and 'h' to express seconds, minutes and hours
func parseDuration(duration string) (int, error) {
	l := len(duration)
	if l < 2 {
		return 0, fmt.Errorf("invalid duation: %s", duration)
	}

	unit := strings.ToLower(duration[l-1:])
	n, err := strconv.ParseInt(duration[:l-1], 10, 32)
	if err != nil {
		return 0, err
	}

	if unit == "s" {
		return int(n), nil
	} else if unit == "m" {
		return int(n) * 60, nil
	} else if unit == "h" {
		return int(n) * 3600, nil
	} else {
		return 0, fmt.Errorf("invalid duration: %s, only supports 's', 'm' and 'h' units", duration)
	}
}
