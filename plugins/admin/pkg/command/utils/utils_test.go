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

package utils

import (
	"testing"

	"gotest.tools/assert"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8sfake "k8s.io/client-go/kubernetes/fake"
)

func TestUtils(t *testing.T) {
	t.Run("report error if ConfigMap not found", func(t *testing.T) {
		desiredCm := &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "config-domain",
				Namespace: "knative-serving",
			},
			Data: make(map[string]string),
		}
		client := k8sfake.NewSimpleClientset()
		err := UpdateConfigMap(client, desiredCm)
		assert.ErrorContains(t, err, "configmaps \"config-domain\" not found", err)
	})

	t.Run("update ConfigMap successfully", func(t *testing.T) {
		oriCm := &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "config-domain",
				Namespace: "knative-serving",
			},
			Data: make(map[string]string),
		}
		desiredCm := &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "config-domain",
				Namespace: "knative-serving",
			},
			Data: map[string]string{
				"dummy.domain": "",
			},
		}
		client := k8sfake.NewSimpleClientset(oriCm)
		err := UpdateConfigMap(client, desiredCm)
		assert.NilError(t, err)

		cm, err := client.CoreV1().ConfigMaps("knative-serving").Get("config-domain", metav1.GetOptions{})
		assert.NilError(t, err)
		assert.Check(t, len(cm.Data) == 1, "expected configmap lengh to be 1")
	})

	t.Run("ConfigMap not changed if desired one is equal to the existed one", func(t *testing.T) {
		oriCm := &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "config-domain",
				Namespace: "knative-serving",
			},
			Data: map[string]string{
				"dummy.domain": "",
			},
		}
		desiredCm := &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "config-domain",
				Namespace: "knative-serving",
			},
			Data: map[string]string{
				"dummy.domain": "",
			},
		}
		client := k8sfake.NewSimpleClientset(oriCm)
		err := UpdateConfigMap(client, desiredCm)
		assert.NilError(t, err)

		updated, err := client.CoreV1().ConfigMaps("knative-serving").Get("config-domain", metav1.GetOptions{})
		assert.NilError(t, err)
		assert.Check(t, equality.Semantic.DeepEqual(updated, oriCm), "configmap should not changed")
	})
}
