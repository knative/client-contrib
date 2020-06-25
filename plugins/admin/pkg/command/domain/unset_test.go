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
	"testing"

	"gotest.tools/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8sfake "k8s.io/client-go/kubernetes/fake"
	"knative.dev/client-contrib/plugins/admin/pkg"

	"knative.dev/client-contrib/plugins/admin/pkg/testutil"
)

func TestNewDomainUnSetCommand(t *testing.T) {

	t.Run("incompleted args", func(t *testing.T) {
		cm := &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      configDomain,
				Namespace: knativeServing,
			},
			Data: make(map[string]string),
		}
		client := k8sfake.NewSimpleClientset(cm)
		p := pkg.AdminParams{
			ClientSet: client,
		}
		cmd := NewDomainUnSetCommand(&p)

		_, err := testutil.ExecuteCommand(cmd, "--custom-domain", "")
		assert.ErrorContains(t, err, "requires the route name", err)
	})

	t.Run("config map not exist", func(t *testing.T) {
		client := k8sfake.NewSimpleClientset()
		p := pkg.AdminParams{
			ClientSet: client,
		}
		cmd := NewDomainUnSetCommand(&p)
		_, err := testutil.ExecuteCommand(cmd, "--custom-domain", "dummy.domain")
		assert.ErrorContains(t, err, "failed to get configmaps", err)
	})

	t.Run("route domain not found", func(t *testing.T) {
		cm := &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      configDomain,
				Namespace: knativeServing,
			},
			Data: map[string]string{
				"dummy.domain": "",
			},
		}
		client := k8sfake.NewSimpleClientset(cm)
		p := pkg.AdminParams{
			ClientSet: client,
		}
		cmd := NewDomainUnSetCommand(&p)

		_, err := testutil.ExecuteCommand(cmd, "--custom-domain", "not-dummy.domain")
		assert.ErrorContains(t, err, "Knative route domain not-dummy.domain not found", err)
	})

	t.Run("unset domain", func(t *testing.T) {
		cm := &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      configDomain,
				Namespace: knativeServing,
			},
			Data: map[string]string{
				"dummy1.domain": "",
				"dummy2.domain": "",
			},
		}
		client := k8sfake.NewSimpleClientset(cm)
		p := pkg.AdminParams{
			ClientSet: client,
		}
		cmd := NewDomainUnSetCommand(&p)
		_, err := testutil.ExecuteCommand(cmd, "--custom-domain", "dummy1.domain")
		assert.NilError(t, err)

		cm, err = client.CoreV1().ConfigMaps(knativeServing).Get(configDomain, metav1.GetOptions{})
		assert.NilError(t, err)
		assert.Check(t, len(cm.Data) == 1, "expected configmap lengh to be 1")

		_, ok := cm.Data["dummy1.domain"]
		assert.Check(t, !ok, "domain key %q should not exists", "dummy1.domain")

		_, ok = cm.Data["dummy2.domain"]
		assert.Check(t, ok, "domain key %q should exists", "dummy2.domain")

		_, err = testutil.ExecuteCommand(cmd, "--custom-domain", "dummy2.domain")
		assert.NilError(t, err)

		cm, err = client.CoreV1().ConfigMaps(knativeServing).Get(configDomain, metav1.GetOptions{})
		assert.NilError(t, err)
		assert.Check(t, len(cm.Data) == 0, "expected configmap lengh to be 0")
	})
}
