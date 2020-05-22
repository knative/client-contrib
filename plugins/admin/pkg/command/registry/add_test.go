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
	"fmt"
	"strings"
	"testing"

	"gotest.tools/assert"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"knative.dev/client-contrib/plugins/admin/pkg"
	"knative.dev/client-contrib/plugins/admin/pkg/testutil"

	k8srand "k8s.io/apimachinery/pkg/util/rand"
	k8sfake "k8s.io/client-go/kubernetes/fake"
	k8stesting "k8s.io/client-go/testing"
)

func TestNewPrAddCommand(t *testing.T) {

	t.Run("incompleted args for registry add", func(t *testing.T) {
		client := k8sfake.NewSimpleClientset()
		p := pkg.AdminParams{
			ClientSet: client,
		}
		cmd := NewPrAddCommand(&p)

		_, err := testutil.ExecuteCommand(cmd, "--username", "")
		assert.ErrorContains(t, err, "requires the registry username")

		_, err = testutil.ExecuteCommand(cmd, "--username", "dummy")
		assert.ErrorContains(t, err, "requires the registry password")

		_, err = testutil.ExecuteCommand(cmd, "--username", "dummy", "--password", "dummy")
		assert.ErrorContains(t, err, "requires the registry server")
	})

	t.Run("missing default service account", func(t *testing.T) {
		client := k8sfake.NewSimpleClientset()
		p := pkg.AdminParams{
			ClientSet: client,
		}

		cmd := NewPrAddCommand(&p)
		_, err := testutil.ExecuteCommand(cmd, "--username", "user", "--password", "dummy", "--server", "docker.io")
		assert.ErrorContains(t, err, "failed to get serviceaccount")
	})

	t.Run("adding registry secret success", func(t *testing.T) {
		sa := corev1.ServiceAccount{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "default",
				Namespace: "default",
			},
		}
		client := k8sfake.NewSimpleClientset(&sa)
		client.PrependReactor("create", "secrets", generateNameReactor)

		p := pkg.AdminParams{
			ClientSet: client,
		}

		cmd := NewPrAddCommand(&p)
		o, err := testutil.ExecuteCommand(cmd, "--username", "user", "--password", "dummy", "--server", "docker.io")
		assert.NilError(t, err)
		assert.Check(t, strings.Contains(o, "Private registry"), "unexpected output: %s", o)

		secrets, err := client.CoreV1().Secrets(sa.Namespace).List(metav1.ListOptions{})
		assert.NilError(t, err)
		assert.Equal(t, len(secrets.Items), 1, "got secrets: %#v", secrets)

		secret := secrets.Items[0]
		assert.Equal(t, secret.Type, corev1.SecretTypeDockerConfigJson)
		assert.Equal(t, secret.GenerateName, "secret-registry-")

		saUpdated, err := client.CoreV1().ServiceAccounts(sa.Namespace).Get(sa.Name, metav1.GetOptions{})
		assert.NilError(t, err)
		assert.Equal(t, len(saUpdated.ImagePullSecrets), 1)
		assert.Equal(t, saUpdated.ImagePullSecrets[0].Name, secret.Name)

		data, ok := secret.Data[".dockerconfigjson"]
		assert.Check(t, ok)

		var r Registry
		err = json.Unmarshal(data, &r)
		assert.NilError(t, err)

		rc, ok := r.Auths["docker.io"]
		assert.Check(t, ok)
		assert.Equal(t, rc.Username, "user")
		assert.Equal(t, rc.Password, "dummy")
		assert.Equal(t, rc.Email, "user@default.email.com")

	})

	t.Run("adding registry secret for service account already have imagepullsecrets", func(t *testing.T) {
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
		client := k8sfake.NewSimpleClientset(&sa)
		client.PrependReactor("create", "secrets", generateNameReactor)

		p := pkg.AdminParams{
			ClientSet: client,
		}

		cmd := NewPrAddCommand(&p)
		o, err := testutil.ExecuteCommand(cmd, "--username", "user", "--password", "dummy", "--server", "docker.io")
		assert.NilError(t, err)
		assert.Check(t, strings.Contains(o, "Private registry"), "unexpected output: %s", o)

		saUpdated, err := client.CoreV1().ServiceAccounts(sa.Namespace).Get(sa.Name, metav1.GetOptions{})
		assert.NilError(t, err)
		assert.Equal(t, len(saUpdated.ImagePullSecrets), 2)
	})
}

func generateNameReactor(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
	s := action.(k8stesting.CreateAction).GetObject().(*corev1.Secret)
	if s.Name == "" && s.GenerateName != "" {
		s.Name = fmt.Sprintf("%s%s", s.GenerateName, k8srand.String(4))
	}
	return false, nil, nil
}
