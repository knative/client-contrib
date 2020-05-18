// Copyright 2020 The Knative Authors

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//     http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package e2e

import (
	"knative.dev/client/lib/test"
)

type E2ETest struct {
	knTest   *test.KnTest
	knPlugin *knPlugin
}

// NewE2ETest for pluginName in pluginPath
func NewE2ETest(pluginName string, pluginPath string, install bool) (*E2ETest, error) {
	knTest, err := test.NewKnTest()
	if err != nil {
		return nil, err
	}

	knPlugin := &knPlugin{
		kn:         knTest.Kn(),
		pluginName: pluginName,
		pluginPath: pluginPath,
		install:    install,
	}

	e2eTest := &E2ETest{
		knTest:   knTest,
		knPlugin: knPlugin,
	}

	return e2eTest, nil
}

// KnTest object
func (e2eTest *E2ETest) KnTest() *test.KnTest {
	return e2eTest.knTest
}

// KnPlugin object
func (e2eTest *E2ETest) KnPlugin() *knPlugin {
	return e2eTest.knPlugin
}
