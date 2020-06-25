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
	"testing"

	"gotest.tools/assert"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8sfake "k8s.io/client-go/kubernetes/fake"
	"knative.dev/client-contrib/plugins/admin/pkg"

	"knative.dev/client-contrib/plugins/admin/pkg/testutil"
)

func TestNewAsUpdateSetCommand(t *testing.T) {
	t.Run("no flags", func(t *testing.T) {
		cm := &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      configAutoscaler,
				Namespace: knativeServing,
			},
			Data: make(map[string]string),
		}
		client := k8sfake.NewSimpleClientset(cm)
		p := pkg.AdminParams{
			ClientSet: client,
		}
		cmd := NewAutoscalingUpdateCommand(&p)

		_, err := testutil.ExecuteCommand(cmd)
		assert.ErrorContains(t, err, "'autoscaling update' requires flag(s)", err)
	})

	t.Run("config map not exist", func(t *testing.T) {
		client := k8sfake.NewSimpleClientset()
		p := pkg.AdminParams{
			ClientSet: client,
		}
		cmd := NewAutoscalingUpdateCommand(&p)
		_, err := testutil.ExecuteCommand(cmd, "--scale-to-zero")
		assert.ErrorContains(t, err, "failed to get ConfigMaps", err)
	})

	t.Run("enable scale-to-zero successfully", func(t *testing.T) {
		cm := &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      configAutoscaler,
				Namespace: knativeServing,
			},
			Data: map[string]string{
				"enable-scale-to-zero": "false",
			},
		}
		client := k8sfake.NewSimpleClientset(cm)
		p := pkg.AdminParams{
			ClientSet: client,
		}
		cmd := NewAutoscalingUpdateCommand(&p)
		_, err := testutil.ExecuteCommand(cmd, "--scale-to-zero")
		assert.NilError(t, err)

		cm, err = client.CoreV1().ConfigMaps(knativeServing).Get(configAutoscaler, metav1.GetOptions{})
		assert.NilError(t, err)
		v, ok := cm.Data["enable-scale-to-zero"]
		assert.Check(t, ok, "key %q should exists", "enable-scale-to-zero")
		assert.Equal(t, "true", v, "enable-scale-to-zero should be true")
	})

	t.Run("disable scale-to-zero successfully", func(t *testing.T) {
		cm := &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      configAutoscaler,
				Namespace: knativeServing,
			},
			Data: map[string]string{
				"enable-scale-to-zero": "true",
			},
		}
		client := k8sfake.NewSimpleClientset(cm)
		p := pkg.AdminParams{
			ClientSet: client,
		}
		cmd := NewAutoscalingUpdateCommand(&p)
		_, err := testutil.ExecuteCommand(cmd, "--no-scale-to-zero")
		assert.NilError(t, err)

		cm, err = client.CoreV1().ConfigMaps(knativeServing).Get(configAutoscaler, metav1.GetOptions{})
		assert.NilError(t, err)
		v, ok := cm.Data["enable-scale-to-zero"]
		assert.Check(t, ok, "key %q should exists", "enable-scale-to-zero")
		assert.Equal(t, "false", v, "enable-scale-to-zero should be false")
	})

	t.Run("enable scale-to-zero but it's already enabled", func(t *testing.T) {
		cm := &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      configAutoscaler,
				Namespace: knativeServing,
			},
			Data: map[string]string{
				"enable-scale-to-zero": "true",
			},
		}
		client := k8sfake.NewSimpleClientset(cm)
		p := pkg.AdminParams{
			ClientSet: client,
		}
		cmd := NewAutoscalingUpdateCommand(&p)

		_, err := testutil.ExecuteCommand(cmd, "--scale-to-zero")
		assert.NilError(t, err)

		updated, err := client.CoreV1().ConfigMaps(knativeServing).Get(configAutoscaler, metav1.GetOptions{})
		assert.NilError(t, err)
		assert.Check(t, equality.Semantic.DeepEqual(updated, cm), "configmap should not be changed")

	})
}
