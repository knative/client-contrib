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
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"gotest.tools/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8srt "k8s.io/apimachinery/pkg/runtime"
	k8sfake "k8s.io/client-go/kubernetes/fake"
	k8sfakecorev1 "k8s.io/client-go/kubernetes/typed/core/v1/fake"
	k8stesting "k8s.io/client-go/testing"
	"knative.dev/client-contrib/plugins/admin/pkg"
	"knative.dev/client-contrib/plugins/admin/pkg/testutil"
)

func newProfilingCommand() *cobra.Command {
	client := k8sfake.NewSimpleClientset(&corev1.ConfigMap{})
	p := pkg.AdminParams{ClientSet: client}
	return NewProfilingCommand(&p)
}

func newProfilingCommandWith(cm *corev1.ConfigMap) (*cobra.Command, *k8sfake.Clientset) {
	client := k8sfake.NewSimpleClientset(cm)
	p := pkg.AdminParams{ClientSet: client}
	return NewProfilingCommand(&p), client
}

type fakeDownloader struct {
	error error
}

func (d *fakeDownloader) Download(t ProfileType, output io.Writer, options ...DownloadOptions) error {
	return d.error
}

func fakeDownloaderBuilder(newError, downloadError error) func(RestConfigGetter, string, string, <-chan struct{}) (ProfileDownloader, error) {
	return func(RestConfigGetter, string, string, <-chan struct{}) (ProfileDownloader, error) {
		return &fakeDownloader{error: downloadError}, newError
	}
}

func removeProfileDataFiles(nameFilter string) {
	files, err := filepath.Glob(nameFilter)
	if err == nil {
		for _, f := range files {
			os.Remove(f)
		}
	}
}

// TestNewProfilingCommand tests the profiling command
func TestNewProfilingCommand(t *testing.T) {
	t.Run("runs profiling without args", func(t *testing.T) {
		out, err := testutil.ExecuteCommand(newProfilingCommand(), "", "")
		assert.NilError(t, err)
		assert.Check(t, strings.Contains(out, "  profiling [flags]"), "expected profiling help output")
	})

	t.Run("runs profiling with conflict args", func(t *testing.T) {
		// --enable and --disable can't be used together
		_, err := testutil.ExecuteCommand(newProfilingCommand(), "--enable", "--disable")
		assert.ErrorContains(t, err, "flags '--enable' and '--disable' can not be used together", err)

		// --enable or --disable --target can't be used with other flags
		argsList := [][]string{
			{"--enable", "--target", "activator"},
			{"--disable", "--target", "activator"},
			{"--enable", "--save-to", "/tmp"},
			{"--disable", "--save-to", "/tmp"},
			{"--enable", "--all"},
			{"--disable", "--all"},
			{"--enable", "--cpu", "5"},
			{"--disable", "--cpu", "5"},
			{"--enable", "--heap"},
			{"--disable", "--heap"},
			{"--enable", "--block"},
			{"--disable", "--block"},
			{"--enable", "--trace", "1m"},
			{"--disable", "--trace", "1m"},
			{"--enable", "--mem-allocs"},
			{"--disable", "--mem-allocs"},
			{"--enable", "--mutex"},
			{"--disable", "--mutex"},
			{"--enable", "--goroutine"},
			{"--disable", "--goroutine"},
			{"--enable", "--thread-create"},
			{"--disable", "--thread-create"},
		}
		for _, args := range argsList {
			_, err := testutil.ExecuteCommand(newProfilingCommand(), args...)
			assert.ErrorContains(t, err, "flag '--enable' or '--disable' can not be used with other flags", err)
		}

		// requires target
		argsList = [][]string{
			{"--save-to", "/tmp"},
			{"--all"},
			{"--all", "--save-to", "/tmp"},
			{"--cpu", "5"},
			{"--cpu", "5", "--save-to", "/tmp"},
			{"--heap"},
			{"--heap", "--save-to", "/tmp"},
			{"--block"},
			{"--block", "--save-to", "/tmp"},
			{"--trace", "1m"},
			{"--trace", "1m", "--save-to", "/tmp"},
			{"--mem-allocs"},
			{"--mem-allocs", "--save-to", "/tmp"},
			{"--mutex"},
			{"--mutex", "--save-to", "/tmp"},
			{"--goroutine"},
			{"--goroutine", "--save-to", "/tmp"},
			{"--thread-create"},
			{"--thread-create", "--save-to", "/tmp"},
		}
		for _, args := range argsList {
			_, err := testutil.ExecuteCommand(newProfilingCommand(), args...)
			assert.ErrorContains(t, err, "requires '--target' flag", err)
		}

		// requires profile type
		argsList = [][]string{
			{"--target", "activator"},
			{"--target", "activator", "--save-to", "/tmp"},
		}
		for _, args := range argsList {
			_, err := testutil.ExecuteCommand(newProfilingCommand(), args...)
			assert.ErrorContains(t, err, "requires '--all' or a specific profiling type flag", err)
		}
	})

	t.Run("parses invalid duration", func(t *testing.T) {
		// invalid numberic
		_, err := parseDuration("1x2")
		assert.ErrorContains(t, err, `parsing "1x2": invalid syntax`, err)

		// invalid unit
		_, err = parseDuration("12x")
		assert.ErrorContains(t, err, "invalid duration: 12x, only supports 's', 'm' and 'h' units", err)

		// not integer
		_, err = parseDuration("1.2")
		assert.ErrorContains(t, err, `parsing "1.2": invalid syntax`, err)
	})

	t.Run("parses valid duration", func(t *testing.T) {
		// empty duration, return default 5
		n, err := parseDuration("")
		assert.NilError(t, err)
		assert.Equal(t, defaultDuration, n)
		n, err = parseDuration("  ")
		assert.NilError(t, err)
		assert.Equal(t, defaultDuration, n)

		// duration is numberic
		n, err = parseDuration("123")
		assert.NilError(t, err)
		assert.Equal(t, 123, n)

		// duration is seconds
		n, err = parseDuration("20s")
		assert.NilError(t, err)
		assert.Equal(t, 20, n)
		n, err = parseDuration("20S")
		assert.NilError(t, err)
		assert.Equal(t, 20, n)

		// duration is minutes
		n, err = parseDuration("2m")
		assert.NilError(t, err)
		assert.Equal(t, 120, n)
		n, err = parseDuration("2M")
		assert.NilError(t, err)
		assert.Equal(t, 120, n)

		// duration is hours
		n, err = parseDuration("1h")
		assert.NilError(t, err)
		assert.Equal(t, 3600, n)
		n, err = parseDuration("1H")
		assert.NilError(t, err)
		assert.Equal(t, 3600, n)
	})

	t.Run("describes duration", func(t *testing.T) {
		assert.Equal(t, "0s", durationDescription(0))
		assert.Equal(t, "-1s", durationDescription(-1))
		assert.Equal(t, "13s", durationDescription(13))
		assert.Equal(t, "1m", durationDescription(60))
		assert.Equal(t, "1m1s", durationDescription(61))
		assert.Equal(t, "1h", durationDescription(3600))
		assert.Equal(t, "1h1s", durationDescription(3601))
		assert.Equal(t, "1h1m", durationDescription(3660))
		assert.Equal(t, "1h1m2s", durationDescription(3662))
	})

	t.Run("failed to enable profiling", func(t *testing.T) {
		cm := &corev1.ConfigMap{}
		cmd, client := newProfilingCommandWith(cm)
		client.CoreV1().(*k8sfakecorev1.FakeCoreV1).PrependReactor("get", "configmaps",
			func(action k8stesting.Action) (handled bool, ret k8srt.Object, err error) {
				return true, &corev1.ConfigMap{}, errors.New("error getting configmap")
			})
		_, err := testutil.ExecuteCommand(cmd, "--enable")
		client.CoreV1().(*k8sfakecorev1.FakeCoreV1).PrependReactor("get", "configmaps",
			func(action k8stesting.Action) (handled bool, ret k8srt.Object, err error) {
				return true, &corev1.ConfigMap{}, errors.New("error getting configmap")
			})
		assert.ErrorContains(t, err, "error getting configmap", err)
	})

	t.Run("failed to disable profiling", func(t *testing.T) {
		cm := &corev1.ConfigMap{}
		cmd, client := newProfilingCommandWith(cm)
		client.CoreV1().(*k8sfakecorev1.FakeCoreV1).PrependReactor("get", "configmaps",
			func(action k8stesting.Action) (handled bool, ret k8srt.Object, err error) {
				return true, &corev1.ConfigMap{}, errors.New("error getting configmap")
			})
		_, err := testutil.ExecuteCommand(cmd, "--disable")
		assert.ErrorContains(t, err, "error getting configmap", err)
	})

	t.Run("successfully enabled profiling", func(t *testing.T) {
		cm := &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{Name: obsConfigMap, Namespace: knNamespace},
			Data:       map[string]string{"profiling.enable": "false"},
		}
		cmd, client := newProfilingCommandWith(cm)
		out, err := testutil.ExecuteCommand(cmd, "--enable")
		assert.NilError(t, err)
		assert.Check(t, strings.Contains(out, "Knative Serving profiling is enabled"))

		newCm, err := client.CoreV1().ConfigMaps(knNamespace).Get(obsConfigMap, metav1.GetOptions{})
		assert.NilError(t, err)
		assert.Equal(t, "true", newCm.Data["profiling.enable"])
	})

	t.Run("successfully disabled profiling", func(t *testing.T) {
		cm := &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{Name: obsConfigMap, Namespace: knNamespace},
			Data:       map[string]string{"profiling.enable": "true"},
		}
		cmd, client := newProfilingCommandWith(cm)
		out, err := testutil.ExecuteCommand(cmd, "--disable")
		assert.NilError(t, err)
		assert.Check(t, strings.Contains(out, "Knative Serving profiling is disabled"))

		newCm, err := client.CoreV1().ConfigMaps(knNamespace).Get(obsConfigMap, metav1.GetOptions{})
		assert.NilError(t, err)
		assert.Equal(t, "false", newCm.Data["profiling.enable"])
	})

	t.Run("save path folder does not exist", func(t *testing.T) {
		cm := &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{Name: obsConfigMap, Namespace: knNamespace},
			Data:       map[string]string{"profiling.enable": "true"},
		}
		cmd, _ := newProfilingCommandWith(cm)
		savePath := "/tmp/xsidsk2hsdks"

		_, err := testutil.ExecuteCommand(cmd, "--target", "activator", "--heap", "--save-to", savePath)
		assert.ErrorContains(t, err, fmt.Sprintf("the specified save path '%s' doesn't exist", savePath), err)
	})

	t.Run("save path is not a folder", func(t *testing.T) {
		cm := &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{Name: obsConfigMap, Namespace: knNamespace},
			Data:       map[string]string{"profiling.enable": "true"},
		}
		cmd, _ := newProfilingCommandWith(cm)
		_, filename, _, _ := runtime.Caller(0)
		savePath := filename

		_, err := testutil.ExecuteCommand(cmd, "--target", "activator", "--heap", "--save-to", savePath)
		assert.ErrorContains(t, err, fmt.Sprintf("the specified save path '%s' is not a folder", savePath), err)
	})

	t.Run("profiling is not enabled when download data", func(t *testing.T) {
		cm := &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{Name: obsConfigMap, Namespace: knNamespace},
			Data:       map[string]string{"profiling.enable": "false"},
		}
		cmd, _ := newProfilingCommandWith(cm)
		_, err := testutil.ExecuteCommand(cmd, "--target", "activator", "--heap")
		assert.ErrorContains(t, err, "profiling is not enabled, please use '--enable' to enalbe it first", err)
	})

	t.Run("failed to get target", func(t *testing.T) {
		cm := &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{Name: obsConfigMap, Namespace: knNamespace},
			Data:       map[string]string{"profiling.enable": "true"},
		}
		cmd, client := newProfilingCommandWith(cm)
		client.CoreV1().(*k8sfakecorev1.FakeCoreV1).PrependReactor("list", "pods",
			func(action k8stesting.Action) (handled bool, ret k8srt.Object, err error) {
				return true, &corev1.PodList{}, errors.New("error listing pods")
			})
		_, err := testutil.ExecuteCommand(cmd, "--target", "activator", "--heap")
		assert.ErrorContains(t, err, "error listing pods", err)
	})

	t.Run("target is not found", func(t *testing.T) {
		cm := &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{Name: obsConfigMap, Namespace: knNamespace},
			Data:       map[string]string{"profiling.enable": "true"},
		}
		cmd, client := newProfilingCommandWith(cm)
		client.CoreV1().(*k8sfakecorev1.FakeCoreV1).PrependReactor("list", "pods",
			func(action k8stesting.Action) (handled bool, ret k8srt.Object, err error) {
				return true, &corev1.PodList{}, nil
			})
		_, err := testutil.ExecuteCommand(cmd, "--target", "activator", "--heap")
		assert.ErrorContains(t, err, "fail to get profiling target 'activator'", err)
	})

	t.Run("failed to get downloader", func(t *testing.T) {
		cm := &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{Name: obsConfigMap, Namespace: knNamespace},
			Data:       map[string]string{"profiling.enable": "true"},
		}
		cmd, client := newProfilingCommandWith(cm)
		pods := corev1.PodList{Items: []corev1.Pod{
			{
				Status: corev1.PodStatus{Phase: corev1.PodRunning},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "activator-1",
					Namespace: knNamespace,
					Labels: map[string]string{
						"app": "activator"},
				},
			},
		}}
		client.CoreV1().(*k8sfakecorev1.FakeCoreV1).PrependReactor("list", "pods",
			func(action k8stesting.Action) (handled bool, ret k8srt.Object, err error) {
				return true, &pods, nil
			})
		newDownloaderFunc = fakeDownloaderBuilder(errors.New("error creating downloader"), nil)
		defer func() { newDownloaderFunc = NewDownloader }()

		_, err := testutil.ExecuteCommand(cmd, "--target", "activator", "--heap")
		assert.ErrorContains(t, err, "error creating downloader", err)
	})

	t.Run("failed to download data", func(t *testing.T) {
		cm := &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{Name: obsConfigMap, Namespace: knNamespace},
			Data:       map[string]string{"profiling.enable": "true"},
		}
		cmd, client := newProfilingCommandWith(cm)
		cwd, _ := os.Getwd()
		podName := "activator-0xxxx"
		pods := corev1.PodList{Items: []corev1.Pod{
			{
				Status: corev1.PodStatus{Phase: corev1.PodRunning},
				ObjectMeta: metav1.ObjectMeta{
					Name:      podName,
					Namespace: knNamespace,
					Labels: map[string]string{
						"app": "activator"},
				},
			},
		}}
		client.CoreV1().(*k8sfakecorev1.FakeCoreV1).PrependReactor("list", "pods",
			func(action k8stesting.Action) (handled bool, ret k8srt.Object, err error) {
				return true, &pods, nil
			})
		newDownloaderFunc = fakeDownloaderBuilder(nil, errors.New("error downloading data"))
		defer func() {
			newDownloaderFunc = NewDownloader
			removeProfileDataFiles(filepath.Join(cwd, podName+"_*"))
		}()

		_, err := testutil.ExecuteCommand(cmd, "--target", podName, "--heap")
		assert.ErrorContains(t, err, "error downloading data", err)
	})

	t.Run("successfully downloaded profiling data for a specific pod target", func(t *testing.T) {
		cm := &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{Name: obsConfigMap, Namespace: knNamespace},
			Data:       map[string]string{"profiling.enable": "true"},
		}
		cmd, client := newProfilingCommandWith(cm)
		cwd, _ := os.Getwd()
		podName := "activator-1xxx"
		pods := corev1.PodList{Items: []corev1.Pod{
			{
				Status: corev1.PodStatus{Phase: corev1.PodRunning},
				ObjectMeta: metav1.ObjectMeta{
					Name:      podName,
					Namespace: knNamespace,
					Labels: map[string]string{
						"app": "activator"},
				},
			},
		}}
		client.CoreV1().(*k8sfakecorev1.FakeCoreV1).PrependReactor("list", "pods",
			func(action k8stesting.Action) (handled bool, ret k8srt.Object, err error) {
				return true, &pods, nil
			})
		newDownloaderFunc = fakeDownloaderBuilder(nil, nil)
		defer func() {
			newDownloaderFunc = NewDownloader
			removeProfileDataFiles(filepath.Join(cwd, podName+"_*"))
		}()

		out, err := testutil.ExecuteCommand(cmd, "--target", podName, "--cpu", "10")
		assert.NilError(t, err)
		expectedMsg := fmt.Sprintf("Saving 10 second(s) cpu profiling data to %s_cpu", filepath.Join(cwd, podName))
		assert.Check(t, strings.Contains(out, expectedMsg), "expected saving cpu profiling data output for"+podName)
	})

	t.Run("successfully downloaded profiling data for a knative component", func(t *testing.T) {
		cm := &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{Name: obsConfigMap, Namespace: knNamespace},
			Data:       map[string]string{"profiling.enable": "true"},
		}
		cmd, client := newProfilingCommandWith(cm)
		cwd, _ := os.Getwd()
		podNames := []string{"activator-2xxx0", "activator-2xxx1"}
		pods := corev1.PodList{}
		for _, name := range podNames {
			pods.Items = append(pods.Items, corev1.Pod{
				Status: corev1.PodStatus{Phase: corev1.PodRunning},
				ObjectMeta: metav1.ObjectMeta{
					Name:      name,
					Namespace: knNamespace,
					Labels: map[string]string{
						"app": "activator"},
				},
			})
		}
		client.CoreV1().(*k8sfakecorev1.FakeCoreV1).PrependReactor("list", "pods",
			func(action k8stesting.Action) (handled bool, ret k8srt.Object, err error) {
				return true, &pods, nil
			})
		newDownloaderFunc = fakeDownloaderBuilder(nil, nil)
		defer func() {
			newDownloaderFunc = NewDownloader
			for _, name := range podNames {
				removeProfileDataFiles(filepath.Join(cwd, name+"_*"))
			}
		}()

		out, err := testutil.ExecuteCommand(cmd, "--target", "activator", "--heap")
		assert.NilError(t, err)
		for _, name := range podNames {
			expectedMsg := fmt.Sprintf("Saving heap profiling data to %s_heap", filepath.Join(cwd, name))
			assert.Check(t, strings.Contains(out, expectedMsg), "expected saving heap profiling data output for "+name)
		}
	})

	t.Run("successfully downloaded multiple profiling types data", func(t *testing.T) {
		cm := &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{Name: obsConfigMap, Namespace: knNamespace},
			Data:       map[string]string{"profiling.enable": "true"},
		}
		cmd, client := newProfilingCommandWith(cm)
		cwd, _ := os.Getwd()
		podName := "activator-3xxxx"
		pods := corev1.PodList{Items: []corev1.Pod{
			{
				Status: corev1.PodStatus{Phase: corev1.PodRunning},
				ObjectMeta: metav1.ObjectMeta{
					Name:      podName,
					Namespace: knNamespace,
					Labels: map[string]string{
						"app": "activator"},
				},
			},
		}}
		client.CoreV1().(*k8sfakecorev1.FakeCoreV1).PrependReactor("list", "pods",
			func(action k8stesting.Action) (handled bool, ret k8srt.Object, err error) {
				return true, &pods, nil
			})
		newDownloaderFunc = fakeDownloaderBuilder(nil, nil)
		defer func() {
			newDownloaderFunc = NewDownloader
			removeProfileDataFiles(filepath.Join(cwd, podName+"_*"))
		}()

		out, err := testutil.ExecuteCommand(cmd, "--target", podName, "--cpu", "8s", "--block", "--mutex")
		assert.NilError(t, err)
		saveFile := filepath.Join(cwd, podName)
		assert.Check(t, strings.Contains(out, "Saving 8 second(s) cpu profiling data to "+saveFile+"_cpu"), "expected saving cpu profiling data output for "+podName)
		assert.Check(t, strings.Contains(out, "Saving block profiling data to "+saveFile+"_block"), "expected saving block profiling data output for "+podName)
		assert.Check(t, strings.Contains(out, "Saving mutex profiling data to "+saveFile+"_mutex"), "expected saving mutex profiling data output for "+podName)
	})

	t.Run("successfully downloaded all profiling types data", func(t *testing.T) {
		cm := &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{Name: obsConfigMap, Namespace: knNamespace},
			Data:       map[string]string{"profiling.enable": "true"},
		}
		cmd, client := newProfilingCommandWith(cm)
		cwd, _ := os.Getwd()
		podName := "activator-4xxxx"
		pods := corev1.PodList{Items: []corev1.Pod{
			{
				Status: corev1.PodStatus{Phase: corev1.PodRunning},
				ObjectMeta: metav1.ObjectMeta{
					Name:      podName,
					Namespace: knNamespace,
					Labels: map[string]string{
						"app": "activator"},
				},
			},
		}}
		client.CoreV1().(*k8sfakecorev1.FakeCoreV1).PrependReactor("list", "pods",
			func(action k8stesting.Action) (handled bool, ret k8srt.Object, err error) {
				return true, &pods, nil
			})
		newDownloaderFunc = fakeDownloaderBuilder(nil, nil)
		defer func() {
			newDownloaderFunc = NewDownloader
			removeProfileDataFiles(filepath.Join(cwd, podName+"_*"))
		}()

		out, err := testutil.ExecuteCommand(cmd, "--target", podName, "--all")
		assert.NilError(t, err)
		saveFile := filepath.Join(cwd, podName)
		expectedMsgs := map[string]string{
			"cpu":           fmt.Sprintf("Saving %d second(s) cpu profiling data to %s_cpu", defaultDuration, saveFile),
			"heap":          fmt.Sprintf("Saving heap profiling data to %s_heap", saveFile),
			"block":         fmt.Sprintf("Saving block profiling data to %s_block", saveFile),
			"trace":         fmt.Sprintf("Saving %d second(s) trace profiling data to %s_trace", defaultDuration, saveFile),
			"mem-allocs":    fmt.Sprintf("Saving mem-allocs profiling data to %s_mem-allocs", saveFile),
			"mutex":         fmt.Sprintf("Saving mutex profiling data to %s_mutex", saveFile),
			"goroutine":     fmt.Sprintf("Saving goroutine profiling data to %s_goroutine", saveFile),
			"thread-create": fmt.Sprintf("Saving thread-create profiling data to %s_thread-create", saveFile),
		}
		for k, v := range expectedMsgs {
			assert.Check(t, strings.Contains(out, v), fmt.Sprintf("expected saving %s profiling data output for %s", k, podName))
		}
	})
}
