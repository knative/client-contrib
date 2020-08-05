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
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"knative.dev/client-contrib/plugins/admin/pkg"
	"knative.dev/client-contrib/plugins/admin/pkg/command/utils"
)

const (
	profilingExample = `
  # To enable Knative Serving profiling
  kn admin profiling --enable

  # To download heap profiling data of autoscaler
  kn admin profiling --target autoscaler --heap

  # To download 2 minutes execution trace data of networking-istio
  kn admin profiling --target networking-istio --trace 2m

  # To download go routing block and memory allocations data of activator and save them to /tmp
  kn admin profiling --target activator --block --mem-allocs --save-to /tmp

  # To download all available profiling data for specified pod activator-5979f56548
  kn admin profiling --target activator-5979f56548 --all
`

	targetFlagUsgae       = "The profiling target. It can be a Knative Serving component name or a specific pod name, e.g: 'activator' or 'activator-586d468c99-w59cm'"
	saveToFlagUsage       = "The path to save the downloaded profiling data, if not speicifed, the data will be saved in current working folder"
	allFlagUsage          = "Download all available profiling data"
	cpuFlagUsage          = "Download cpu profiling data, you can specify a profiling data duration with 's' for second(s), 'm' for minute(s) and 'h' for hour(s), e.g: '1m' for one minute"
	heapFlagUsage         = "Download heap profiling data"
	blockFlagUsage        = "Download go routine blocking data"
	traceFlagUsage        = "Download execution trace data, you can specify a trace data duration with 's' for second(s), 'm' for minute(s) and 'h' for hour(s), e.g: '1m' for one minute"
	memAllocsFlagUsage    = "Download memory allocations data"
	mutexFlagUsage        = "Download holders of contended mutexes data"
	goroutineFlagUsage    = "Download stack traces of all current goroutines data"
	threadCreateFlagUsage = "Download stack traces that led to the creation of new OS threads data"

	cpuFlagName          = "cpu"
	heapFlagName         = "heap"
	blockFlagName        = "block"
	traceFlagName        = "trace"
	memAllocsFlagName    = "mem-allocs"
	mutexFlagName        = "mutex"
	goroutineFlagName    = "goroutine"
	threadCreateFlagName = "thread-create"
	knNamespace          = "knative-serving"
	obsConfigMap         = "config-observability"
	defaultDuration      = 5
	defaultProfilingTime = OptionProfilingTime(defaultDuration * time.Second)
)

// profilingFlags defines flag values for profiling command
type profilingFlags struct {
	enable              bool
	disable             bool
	target              string
	saveTo              string
	allProfiles         bool
	cpuProfile          string
	heapProfile         bool
	blockProfile        bool
	traceProfile        string
	memAllocsProfile    bool
	mutexProfile        bool
	goroutineProfile    bool
	threadCreateProfile bool
}

// profileTypeOption is a helper struct to download profile type data
type profileTypeOption struct {
	profileType    ProfileType
	downloadOption DownloadOptions
}

// declare newDownloaderFunc variable to help us write UT
var newDownloaderFunc = NewDownloader

// NewProfilingCommand creates a profiling command
func NewProfilingCommand(p *pkg.AdminParams) *cobra.Command {
	pflags := profilingFlags{}

	var profilingCmd = &cobra.Command{
		Use:     "profiling",
		Aliases: []string{"prof"},
		Short:   "Profiling Knative Serving components",
		Long:    `Enable Knative Serving components profiling and download profiling data`,
		Example: profilingExample,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			flags := cmd.Flags()
			// no flag given, print help
			if flags.NFlag() < 1 {
				cmd.Help()
				return nil
			}

			isEnableSet := flags.Changed("enable")
			isDisableSet := flags.Changed("disable")
			isTargetSet := flags.Changed("target")
			isSaveToSet := flags.Changed("save-to")
			isAllProfilesSet := flags.Changed("all")
			isProfileTypeSet := (flags.Changed(cpuFlagName) || flags.Changed(heapFlagName) || flags.Changed(blockFlagName) ||
				flags.Changed(traceFlagName) || flags.Changed(memAllocsFlagName) || flags.Changed(mutexFlagName) ||
				flags.Changed(goroutineFlagName) || flags.Changed(threadCreateFlagName))

			// enable and disable can't be used togerther
			if isEnableSet && isDisableSet {
				return fmt.Errorf("flags '--enable' and '--disable' can not be used together")
			}

			// enable or disable can't be used with other flags
			if (isEnableSet || isDisableSet) && (isTargetSet || isProfileTypeSet || isAllProfilesSet || isSaveToSet) {
				return fmt.Errorf("flag '--enable' or '--disable' can not be used with other flags")
			}

			// --target flag is needed
			if !isTargetSet && (isProfileTypeSet || isAllProfilesSet || isSaveToSet) {
				return fmt.Errorf("requires '--target' flag")
			}

			// --profile-type is needed
			if !isProfileTypeSet && !isAllProfilesSet && (isTargetSet || isSaveToSet) {
				return fmt.Errorf("requires '--all' or a specific profiling type flag")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			flags := cmd.Flags()
			if flags.NFlag() < 1 {
				return nil
			} else if flags.Changed("enable") {
				return configProfiling(p.ClientSet, cmd, true)
			} else if flags.Changed("disable") {
				return configProfiling(p.ClientSet, cmd, false)
			} else {
				return downloadProfileData(p, cmd, &pflags)
			}
		},
	}

	flags := profilingCmd.Flags()
	flags.BoolVar(&pflags.enable, "enable", false, "Enable Knative Serving profiling")
	flags.BoolVar(&pflags.disable, "disable", false, "Disable Knative Serving profiling")
	flags.StringVarP(&pflags.target, "target", "t", "", targetFlagUsgae)
	flags.StringVarP(&pflags.saveTo, "save-to", "s", "", saveToFlagUsage)
	flags.BoolVar(&pflags.allProfiles, "all", false, allFlagUsage)
	flags.StringVarP(&pflags.cpuProfile, cpuFlagName, "", "5s", cpuFlagUsage)
	flags.BoolVar(&pflags.heapProfile, heapFlagName, false, heapFlagUsage)
	flags.BoolVar(&pflags.blockProfile, blockFlagName, false, blockFlagUsage)
	flags.StringVarP(&pflags.traceProfile, traceFlagName, "", "5s", traceFlagUsage)
	flags.BoolVar(&pflags.memAllocsProfile, memAllocsFlagName, false, memAllocsFlagUsage)
	flags.BoolVar(&pflags.mutexProfile, mutexFlagName, false, mutexFlagUsage)
	flags.BoolVar(&pflags.goroutineProfile, goroutineFlagName, false, goroutineFlagUsage)
	flags.BoolVar(&pflags.threadCreateProfile, threadCreateFlagName, false, threadCreateFlagUsage)
	return profilingCmd
}

// configProfiling enables or disables knative profiling
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

	err = utils.UpdateConfigMap(c, desiredCm)
	if err != nil {
		return fmt.Errorf("failed to update ConfigMap %s in namespace %s: %+v", obsConfigMap, knNamespace, err)
	}

	if enable {
		cmd.Println("Knative Serving profiling is enabled")
	} else {
		cmd.Println("Knative Serving profiling is disabled")
	}
	return nil
}

// isProfilingEnabled checks if the profiling is enabled
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
	duration = strings.TrimSpace(duration)
	if duration == "" {
		return defaultDuration, nil
	}

	l := len(duration)
	unit := duration[l-1]
	// duration is numberic
	if unit >= '0' && unit <= '9' {
		n, err := strconv.ParseInt(duration, 10, 32)
		if err != nil {
			return 0, err
		}
		return int(n), nil
	}

	// parse duration[:l-1] as int
	n, err := strconv.ParseInt(duration[:l-1], 10, 32)
	if err != nil {
		return 0, err
	}

	// calculate duration by unit
	if unit == 's' || unit == 'S' {
		return int(n), nil
	} else if unit == 'm' || unit == 'M' {
		return int(n) * 60, nil
	} else if unit == 'h' || unit == 'H' {
		return int(n) * 3600, nil
	} else {
		return 0, fmt.Errorf("invalid duration: %s, only supports 's', 'm' and 'h' units", duration)
	}
}

// durationDescription describes a given seconds as 0h0m0s format
func durationDescription(seconds int64) string {
	if seconds < 1 {
		return fmt.Sprintf("%ds", seconds)
	}

	s := ""
	if hours := seconds / 3600; hours > 0 {
		s = fmt.Sprintf("%dh", hours)
	}
	left := seconds % 3600
	if mins := left / 60; mins > 0 {
		s = fmt.Sprintf("%s%dm", s, mins)
	}
	if secs := left % 60; secs > 0 {
		s = fmt.Sprintf("%s%ds", s, secs)
	}
	return s
}

// downloadProfileData downloads profile data by given profile type
func downloadProfileData(p *pkg.AdminParams, cmd *cobra.Command, pflags *profilingFlags) error {
	// check profile types
	profileTypes := map[string]profileTypeOption{}
	if pflags.allProfiles {
		// all profile types
		profileTypes[cpuFlagName] = profileTypeOption{
			profileType:    ProfileTypeProfile,
			downloadOption: defaultProfilingTime,
		}
		profileTypes[heapFlagName] = profileTypeOption{profileType: ProfileTypeProfile}
		profileTypes[blockFlagName] = profileTypeOption{profileType: ProfileTypeBlock}
		profileTypes[traceFlagName] = profileTypeOption{
			profileType:    ProfileTypeTrace,
			downloadOption: defaultProfilingTime,
		}
		profileTypes[memAllocsFlagName] = profileTypeOption{profileType: ProfileTypeAllocs}
		profileTypes[mutexFlagName] = profileTypeOption{profileType: ProfileTypeMutex}
		profileTypes[goroutineFlagName] = profileTypeOption{profileType: ProfileTypeGoroutine}
		profileTypes[threadCreateFlagName] = profileTypeOption{profileType: ProfileTypeThreadCreate}
	} else {
		flags := cmd.Flags()
		// cpu profile type
		if flags.Changed(cpuFlagName) {
			op := profileTypeOption{profileType: ProfileTypeProfile}
			if pflags.cpuProfile == "" {
				op.downloadOption = defaultProfilingTime
			} else {
				duration, err := parseDuration(pflags.cpuProfile)
				if err != nil {
					return err
				}
				op.downloadOption = OptionProfilingTime(time.Duration(duration) * time.Second)
			}
			profileTypes[cpuFlagName] = op
		}
		// heap profile type
		if flags.Changed(heapFlagName) {
			profileTypes[heapFlagName] = profileTypeOption{profileType: ProfileTypeHeap}
		}
		// block profile type
		if flags.Changed(blockFlagName) {
			profileTypes[blockFlagName] = profileTypeOption{profileType: ProfileTypeBlock}
		}
		// trace profile type
		if flags.Changed(traceFlagName) {
			op := profileTypeOption{profileType: ProfileTypeTrace}
			if pflags.traceProfile == "" {
				op.downloadOption = defaultProfilingTime
			} else {
				duration, err := parseDuration(pflags.traceProfile)
				if err != nil {
					return err
				}
				op.downloadOption = OptionProfilingTime(time.Duration(duration) * time.Second)
			}
			profileTypes[traceFlagName] = op
		}
		// mem-allocs profile type
		if flags.Changed(memAllocsFlagName) {
			profileTypes[memAllocsFlagName] = profileTypeOption{profileType: ProfileTypeAllocs}
		}
		// mutex profile type
		if flags.Changed(mutexFlagName) {
			profileTypes[mutexFlagName] = profileTypeOption{profileType: ProfileTypeMutex}
		}
		// goroutine profile type
		if flags.Changed(goroutineFlagName) {
			profileTypes[goroutineFlagName] = profileTypeOption{profileType: ProfileTypeGoroutine}
		}
		// thread-create profile type
		if flags.Changed(threadCreateFlagName) {
			profileTypes[threadCreateFlagName] = profileTypeOption{profileType: ProfileTypeThreadCreate}
		}
	}

	// check --save-to path
	var err error
	if pflags.saveTo == "" {
		pflags.saveTo, err = os.Getwd()
		if err != nil {
			return fmt.Errorf("fail to get current working folder to save profiling data: %+v", err)
		}
	} else {
		stat, err := os.Stat(pflags.saveTo)
		if os.IsNotExist(err) {
			return fmt.Errorf("the specified save path '%s' doesn't exist", pflags.saveTo)
		} else if err != nil {
			return err
		}
		if !stat.IsDir() {
			return fmt.Errorf("the specified save path '%s' is not a folder", pflags.saveTo)
		}
	}

	// check if profiling is enabled, if not, print message to ask user enable it first
	enabled, err := isProfilingEnabled(p.ClientSet)
	if err != nil {
		return err
	}
	if !enabled {
		return fmt.Errorf("profiling is not enabled, please use '--enable' to enalbe it first")
	}

	// try to find target as a knative component name
	pods, err := p.ClientSet.CoreV1().Pods(knNamespace).List(metav1.ListOptions{LabelSelector: "app=" + pflags.target})
	if err != nil {
		return err
	}
	// if no pod found, try to find target as a pod name in knative namespace
	if len(pods.Items) < 1 {
		pods, err = p.ClientSet.CoreV1().Pods(knNamespace).List(metav1.ListOptions{})
		if err != nil {
			return err
		}
		if len(pods.Items) < 1 {
			return fmt.Errorf("fail to get profiling target '%s'", pflags.target)
		}

		// check if target is found as pod
		found := false
		for _, p := range pods.Items {
			if p.Name == pflags.target {
				pods.Items = []corev1.Pod{p}
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("fail to get profiling target '%s'", pflags.target)
		}
	}

	// iterates pods to download specified profiling data
	for _, pod := range pods.Items {
		err = func() error {
			cmd.Printf("Starting to download profiling data for pod %s...\n", pod.Name)
			end := make(chan struct{})
			downloader, err := newDownloaderFunc(p, pod.Name, knNamespace, end)
			if err != nil {
				return err
			}
			defer close(end)

			// iterates specified profile types to download data
			for k, v := range profileTypes {
				duration := ""
				filename := pod.Name + "_" + k
				options := []DownloadOptions{}
				if t, ok := v.downloadOption.(OptionProfilingTime); ok {
					seconds := int64(time.Duration(t) / time.Second)
					duration = strconv.FormatInt(seconds, 10) + " second(s) "
					filename += "_" + durationDescription(seconds)
					options = append(options, v.downloadOption)
				}
				filename += "_" + time.Now().Format("20060102150405")
				dataFilePath := filepath.Join(pflags.saveTo, filename)
				f, err := os.Create(dataFilePath)
				if err != nil {
					return err
				}

				cmd.Printf("Saving %s%s profiling data to %s\n", duration, k, dataFilePath)
				err = downloader.Download(v.profileType, f, options...)
				f.Close()
				if err != nil {
					return err
				}
			}
			return nil
		}()
		if err != nil {
			return err
		}
	}
	return nil
}
