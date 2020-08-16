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

package registry

import (
	"encoding/json"
	"strings"
	"testing"

	"gotest.tools/assert"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"knative.dev/client-contrib/plugins/admin/pkg"
	"knative.dev/client-contrib/plugins/admin/pkg/testutil"

	k8sfake "k8s.io/client-go/kubernetes/fake"
)

func TestNewRegistryRmCommand(t *testing.T) {
	t.Run("incompleted args for registry remove", func(t *testing.T) {
		client := k8sfake.NewSimpleClientset()

		p := &pkg.AdminParams{
			ClientSet: client,
		}

		cmd := NewRegistryRmCommand(p)

		_, err := testutil.ExecuteCommand(cmd, "--username", "")
		assert.ErrorContains(t, err, "requires the registry username")

		_, err = testutil.ExecuteCommand(cmd, "--username", "dummy", "--server", "")
		assert.ErrorContains(t, err, "requires the registry server")
	})

	t.Run("registry not found", func(t *testing.T) {
		client := k8sfake.NewSimpleClientset()

		p := &pkg.AdminParams{
			ClientSet: client,
		}
		cmd := NewRegistryRmCommand(p)
		o, err := testutil.ExecuteCommand(cmd, "--username", "user", "--server", "docker.io")
		assert.NilError(t, err)
		assert.Check(t, strings.Contains(o, "No registry found"), "unexpected output: %s", o)
	})

	t.Run("registry removed successfully in default namespace using default serviceaccount", func(t *testing.T) {
		sa := corev1.ServiceAccount{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "default",
				Namespace: "default",
			},
			ImagePullSecrets: []corev1.LocalObjectReference{
				{
					Name: "dummy-secret",
				},
			},
		}

		dockerCfg := Registry{
			Auths: Auths{
				"docker.io": registryCred{
					Username: "user",
					Password: "password",
					Email:    "email",
				},
			},
		}

		j, err := json.Marshal(dockerCfg)
		assert.NilError(t, err)

		secret := corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "dummy-secret",
				Namespace: "default",
				Labels: map[string]string{
					pkg.LabelManagedBy: AdminRegistryCmdName,
				},
			},
			Data: map[string][]byte{
				".dockerconfigjson": j,
			},
		}
		client := k8sfake.NewSimpleClientset(&sa, &secret)

		p := &pkg.AdminParams{
			ClientSet: client,
		}

		cmd := NewRegistryRmCommand(p)
		o, err := testutil.ExecuteCommand(cmd, "--username", "user", "--server", "docker.io")
		assert.NilError(t, err)
		assert.Check(t, strings.Contains(o, "ImagePullSecrets of serviceaccount 'default' in namespace 'default' is updated"), "unexpected output: %s", o)
		assert.Check(t, strings.Contains(o, "Secret 'dummy-secret' in namespace 'default' is deleted"), "unexpected output: %s", o)

		_, err = client.CoreV1().Secrets("default").Get("dummy-secret", metav1.GetOptions{})
		assert.ErrorContains(t, err, "not found")
		saUpdated, err := client.CoreV1().ServiceAccounts("default").Get("default", metav1.GetOptions{})
		assert.NilError(t, err)
		isContain := false
		for _, imagePullSecret := range saUpdated.ImagePullSecrets {
			if imagePullSecret.Name == "dummy-secret" {
				isContain = true
				break
			}
		}
		assert.Check(t, !isContain, "ImagePullSecrets in the updated serviceaccount should not contain the removed secret")
	})

	t.Run("registry removed successfully in custom namespace using custom serviceaccount", func(t *testing.T) {
		ns := corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: "custom-namespace",
			},
		}
		sa := corev1.ServiceAccount{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "custom-serviceaccount",
				Namespace: "custom-namespace",
			},
			ImagePullSecrets: []corev1.LocalObjectReference{
				{
					Name: "dummy-secret",
				},
			},
		}

		dockerCfg := Registry{
			Auths: Auths{
				"docker.io": registryCred{
					Username: "user",
					Password: "password",
					Email:    "email",
				},
			},
		}

		j, err := json.Marshal(dockerCfg)
		assert.NilError(t, err)

		secret := corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "dummy-secret",
				Namespace: "custom-namespace",
				Labels: map[string]string{
					pkg.LabelManagedBy: AdminRegistryCmdName,
				},
			},
			Data: map[string][]byte{
				".dockerconfigjson": j,
			},
		}
		client := k8sfake.NewSimpleClientset(&ns, &sa, &secret)

		p := &pkg.AdminParams{
			ClientSet: client,
		}

		cmd := NewRegistryRmCommand(p)
		o, err := testutil.ExecuteCommand(cmd, "--username", "user", "--server", "docker.io", "--namespace", "custom-namespace", "--serviceaccount", "custom-serviceaccount")
		assert.NilError(t, err)
		assert.Check(t, strings.Contains(o, "ImagePullSecrets of serviceaccount 'custom-serviceaccount' in namespace 'custom-namespace' is updated"), "unexpected output: %s", o)
		assert.Check(t, strings.Contains(o, "Secret 'dummy-secret' in namespace 'custom-namespace' is deleted"), "unexpected output: %s", o)

		_, err = client.CoreV1().Secrets(ns.Name).Get("dummy-secret", metav1.GetOptions{})
		assert.ErrorContains(t, err, "not found")
		saUpdated, err := client.CoreV1().ServiceAccounts(ns.Name).Get("custom-serviceaccount", metav1.GetOptions{})
		assert.NilError(t, err)
		isContain := false
		for _, imagePullSecret := range saUpdated.ImagePullSecrets {
			if imagePullSecret.Name == "dummy-secret" {
				isContain = true
				break
			}
		}
		assert.Check(t, !isContain, "ImagePullSecrets in the updated serviceaccount should not contain the removed secret")
	})

}
