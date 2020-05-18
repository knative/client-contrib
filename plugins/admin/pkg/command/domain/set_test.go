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
	"gopkg.in/yaml.v2"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8sfake "k8s.io/client-go/kubernetes/fake"
	"knative.dev/client-contrib/plugins/admin/pkg"
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
				Name:      "config-domain",
				Namespace: "knative-serving",
			},
			Data: make(map[string]string),
		}
		client := k8sfake.NewSimpleClientset(cm)
		p := pkg.AdminParams{
			ClientSet: client,
		}
		cmd := NewDomainSetCommand(&p)
		_, _, err := executeCommandC(cmd, "--custom-domain", "")
		if err == nil {
			t.Errorf("expected error when config-domain is empty")
		}
	})

	t.Run("config map not exist", func(t *testing.T) {
		client := k8sfake.NewSimpleClientset()
		p := pkg.AdminParams{
			ClientSet: client,
		}
		cmd := NewDomainSetCommand(&p)
		_, _, err := executeCommandC(cmd, "--custom-domain", "dummy.domain")
		if err == nil {
			t.Errorf("expected error when config-domain configmap does not exist")
		}
		if !strings.HasPrefix(err.Error(), "Failed to get ConfigMaps:") {
			t.Errorf("unexpected error string: %+v", err)
		}
	})

	t.Run("setting domain config without selector", func(t *testing.T) {
		cm := &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "config-domain",
				Namespace: "knative-serving",
			},
			Data: make(map[string]string),
		}
		client := k8sfake.NewSimpleClientset(cm)
		p := pkg.AdminParams{
			ClientSet: client,
		}
		cmd := NewDomainSetCommand(&p)
		_, _, err := executeCommandC(cmd, "--custom-domain", "dummy.domain")
		if err != nil {
			t.Errorf("unexpected error %+v", err)
		}

		cm, err = client.CoreV1().ConfigMaps("knative-serving").Get("config-domain", metav1.GetOptions{})
		if err != nil {
			t.Error(err)
		}

		if len(cm.Data) != 1 {
			t.Errorf("expected configmap lengh to be 1, actual %d", len(cm.Data))
		}
		v, ok := cm.Data["dummy.domain"]
		if !ok {
			t.Errorf("domain key %q does not exists", "dummy.domain")
		}
		if v != "" {
			t.Errorf("value of key domain is not empty: %q", v)
		}
	})

	t.Run("setting domain config with unchanged value", func(t *testing.T) {
		cm := &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "config-domain",
				Namespace: "knative-serving",
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
		_, _, err := executeCommandC(cmd, "--custom-domain", "dummy.domain")
		if err != nil {
			t.Errorf("unexpected error %+v", err)
		}

		updated, err := client.CoreV1().ConfigMaps("knative-serving").Get("config-domain", metav1.GetOptions{})
		if err != nil {
			t.Error(err)
		}

		if !equality.Semantic.DeepEqual(updated, cm) {
			t.Error("configmap should not changed")
		}

	})

	t.Run("adding domain config without selector with existing domain configuration", func(t *testing.T) {
		cm := &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "config-domain",
				Namespace: "knative-serving",
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
		_, o, err := executeCommandC(cmd, "--custom-domain", "dummy.domain")
		if err != nil {
			t.Errorf("unexpected error %+v", err)
		}
		if o == "" {
			t.Error("expected update information in standard output. got empty.")
		}

		cm, err = client.CoreV1().ConfigMaps("knative-serving").Get("config-domain", metav1.GetOptions{})
		if err != nil {
			t.Error(err)
		}

		if len(cm.Data) != 1 {
			t.Errorf("expected configmap lengh to be 1, actual %d", len(cm.Data))
		}
		v, ok := cm.Data["dummy.domain"]
		if !ok {
			t.Errorf("domain key %q does not exists", "dummy.domain")
		}
		if v != "" {
			t.Errorf("value of key domain is not empty: %q", v)
		}

	})

	t.Run("adding domain config with selector", func(t *testing.T) {
		cm := &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "config-domain",
				Namespace: "knative-serving",
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
		_, o, err := executeCommandC(cmd, "--custom-domain", "dummy.domain", "--selector", "app=dummy")
		if err != nil {
			t.Errorf("unexpected error %+v", err)
		}
		if o == "" {
			t.Error("expected update information in standard output. got empty.")
		}

		cm, err = client.CoreV1().ConfigMaps("knative-serving").Get("config-domain", metav1.GetOptions{})
		if err != nil {
			t.Error(err)
		}

		if len(cm.Data) != 2 {
			t.Errorf("expected configmap lengh to be 2, actual %d", len(cm.Data))
		}
		v, ok := cm.Data["dummy.domain"]
		if !ok {
			t.Errorf("domain key %q does not exists", "dummy.domain")
		}

		var s domainSelector
		err = yaml.Unmarshal([]byte(v), &s)
		if err != nil {
			t.Errorf("unmarshal domain config error %v", v)
		}
		if len(s.Selector) != 1 {
			t.Errorf("selector should only contain one key-value pair, got %+v", s.Selector)
		}
		v, ok = s.Selector["app"]
		if !ok {
			t.Errorf("key %q dose not exist", "app")
		}
		if v != "dummy" {
			t.Errorf("got unexpected value %q", v)
		}
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
