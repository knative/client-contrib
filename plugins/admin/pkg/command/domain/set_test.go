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
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v2"
	"gotest.tools/assert"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8sfake "k8s.io/client-go/kubernetes/fake"
	"knative.dev/client-contrib/plugins/admin/pkg"

	"knative.dev/client-contrib/plugins/admin/pkg/testutil"
)

type domainSelector struct {
	Selector map[string]string `yaml:"selector,omitempty"`
}

func executeCommandC(root *cobra.Command, args ...string) (c *cobra.Command, output string, err error) {
	buf := new(bytes.Buffer)
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(args)
	c, err = root.ExecuteC()
	return c, buf.String(), err
}

func TestNewDomainSetCommand(t *testing.T) {

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
		cmd := NewDomainSetCommand(&p)

		_, err := testutil.ExecuteCommand(cmd, "--custom-domain", "")
		assert.ErrorContains(t, err, "requires the route name", err)
	})

	t.Run("config map not exist", func(t *testing.T) {
		client := k8sfake.NewSimpleClientset()
		p := pkg.AdminParams{
			ClientSet: client,
		}
		cmd := NewDomainSetCommand(&p)
		_, err := testutil.ExecuteCommand(cmd, "--custom-domain", "dummy.domain")
		assert.ErrorContains(t, err, "failed to get ConfigMap", err)
	})

	t.Run("setting domain config without selector", func(t *testing.T) {
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
		cmd := NewDomainSetCommand(&p)
		_, err := testutil.ExecuteCommand(cmd, "--custom-domain", "dummy.domain")
		assert.NilError(t, err)

		cm, err = client.CoreV1().ConfigMaps(knativeServing).Get(configDomain, metav1.GetOptions{})
		assert.NilError(t, err)
		assert.Check(t, len(cm.Data) == 1, "expected configmap lengh to be 1")

		v, ok := cm.Data["dummy.domain"]
		assert.Check(t, ok, "domain key %q should exists", "dummy.domain")
		assert.Equal(t, "", v, "value of key domain should be empty")
	})

	t.Run("setting domain config with unchanged value", func(t *testing.T) {
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
		cmd := NewDomainSetCommand(&p)

		_, err := testutil.ExecuteCommand(cmd, "--custom-domain", "dummy.domain")
		assert.NilError(t, err)

		updated, err := client.CoreV1().ConfigMaps(knativeServing).Get(configDomain, metav1.GetOptions{})
		assert.NilError(t, err)
		assert.Check(t, equality.Semantic.DeepEqual(updated, cm), "configmap should not changed")

	})

	t.Run("adding domain config without selector with existing domain configuration", func(t *testing.T) {
		cm := &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      configDomain,
				Namespace: knativeServing,
			},
			Data: map[string]string{
				"foo.bar": "",
			},
		}
		client := k8sfake.NewSimpleClientset(cm)
		p := pkg.AdminParams{
			ClientSet: client,
		}
		cmd := NewDomainSetCommand(&p)
		o, err := testutil.ExecuteCommand(cmd, "--custom-domain", "dummy.domain")
		assert.NilError(t, err)
		assert.Check(t, strings.Contains(o, "Set knative route domain \"dummy.domain\""), "expected update information in standard output")

		cm, err = client.CoreV1().ConfigMaps(knativeServing).Get(configDomain, metav1.GetOptions{})
		assert.NilError(t, err)
		assert.Check(t, len(cm.Data) == 1, "expected configmap lengh to be 1, actual %d", len(cm.Data))

		v, ok := cm.Data["dummy.domain"]
		assert.Check(t, ok, "domain key %q should exists", "dummy.domain")
		assert.Equal(t, "", v, "value of key domain should be empty")
	})

	t.Run("adding domain config with selector", func(t *testing.T) {
		cm := &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      configDomain,
				Namespace: knativeServing,
			},
			Data: map[string]string{
				"foo.bar": "",
			},
		}
		client := k8sfake.NewSimpleClientset(cm)
		p := pkg.AdminParams{
			ClientSet: client,
		}
		cmd := NewDomainSetCommand(&p)

		o, err := testutil.ExecuteCommand(cmd, "--custom-domain", "dummy.domain", "--selector", "app=dummy")
		assert.NilError(t, err)
		assert.Check(t, strings.Contains(o, "Set knative route domain \"dummy.domain\" with selector [app=dummy]"), "invalid output %q", o)

		cm, err = client.CoreV1().ConfigMaps(knativeServing).Get(configDomain, metav1.GetOptions{})
		assert.NilError(t, err)
		assert.Check(t, len(cm.Data) == 2, "expected configmap lengh to be 2, actual %d", len(cm.Data))

		v, ok := cm.Data["dummy.domain"]
		assert.Check(t, ok, "domain key %q should exists", "dummy.domain")

		var s domainSelector
		err = yaml.Unmarshal([]byte(v), &s)
		assert.NilError(t, err)
		assert.Check(t, len(s.Selector) == 1, "selector should only contain one key-value pair, got %+v", s.Selector)

		v, ok = s.Selector["app"]
		assert.Check(t, ok, "key %q should exist", "app")
		assert.Equal(t, "dummy", v)
	})

	t.Run("adding domain config with invalid selector", func(t *testing.T) {
		cm := &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      configDomain,
				Namespace: knativeServing,
			},
			Data: map[string]string{
				"foo.bar": "",
			},
		}
		client := k8sfake.NewSimpleClientset(cm)
		p := pkg.AdminParams{
			ClientSet: client,
		}
		cmd := NewDomainSetCommand(&p)

		_, err := testutil.ExecuteCommand(cmd, "--custom-domain", "dummy.domain", "--selector", "app")
		assert.ErrorContains(t, err, "expecting the selector format 'name=value', found 'app'", err)
	})
}

func Test_splitByEqualSign(t *testing.T) {
	tests := []struct {
		name    string
		pair    string
		k       string
		v       string
		wantErr bool
	}{
		{"normal case", "app=abc", "app", "abc", false},
		{"normal case with spaces", " app=abc ", "app", "abc", false},
		{"empty key and value", "=", "", "", true},
		{"space key and value", " = ", "", "", true},
		{"empty key 1", "=abc", "", "", true},
		{"empty key 2", " =abc", "", "", true},
		{"empty value 1", "app=", "", "", true},
		{"empty value 2", "app= ", "", "", true},
		{"invalid input 1", "app=aaa=bbb", "", "", true},
		{"invalid input 2", "app.123", "", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotk, gotv, err := splitByEqualSign(tt.pair)
			if (err != nil) != tt.wantErr {
				t.Errorf("splitByEqualSign() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotk != tt.k {
				t.Errorf("splitByEqualSign() got = %v, want %v", gotk, tt.k)
			}
			if gotv != tt.v {
				t.Errorf("splitByEqualSign() got1 = %v, want %v", gotv, tt.v)
			}
		})
	}
}
