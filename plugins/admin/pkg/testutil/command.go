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

package testutil

import (
	"bytes"

	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	k8sfake "k8s.io/client-go/kubernetes/fake"

	"knative.dev/client-contrib/plugins/admin/pkg"
)

// ExecuteCommandC execute cobra.command and catch the output
func ExecuteCommandC(root *cobra.Command, args ...string) (c *cobra.Command, output string, err error) {
	buf := new(bytes.Buffer)
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(args)
	c, err = root.ExecuteC()
	return c, buf.String(), err
}

// ExecuteCommand similar to ExecuteCommandC but does not return *cobra.Command
func ExecuteCommand(root *cobra.Command, args ...string) (output string, err error) {
	_, o, err := ExecuteCommandC(root, args...)
	return o, err
}

// NewTestAdminParams creates an AdminParams and kubenetes clientset for testing
func NewTestAdminParams(objects ...runtime.Object) (*pkg.AdminParams, *k8sfake.Clientset) {
	client := k8sfake.NewSimpleClientset(objects...)
	return &pkg.AdminParams{
		NewKubeClient: func() (kubernetes.Interface, error) {
			return client, nil
		},
	}, client
}
