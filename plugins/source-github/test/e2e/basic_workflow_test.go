// Copyright 2019 The Knative Authors

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//     http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// +build e2e
// +build !eventing

package e2e

import (
	"os"
	"path/filepath"
	"testing"

	"gotest.tools/assert"

	"knative.dev/client/lib/test"
	"knative.dev/client/pkg/util"

	commone2e "github.com/maximilien/kn-source-pkg/test/e2e"
)

type e2eTest struct {
	it *commone2e.E2ETest
}

func newE2ETest(t *testing.T) *e2eTest {
	currentDir, err := os.Getwd()
	assert.NilError(t, err)

	it, err := commone2e.NewE2ETest("kn-source_github", filepath.Join(currentDir, "../.."), false)
	assert.NilError(t, err)
	defer func() {
		assert.NilError(t, it.KnTest().Teardown())
	}()

	e2eTest := &e2eTest{
		it: it,
	}

	return e2eTest
}

func TestBasicWorkflow(t *testing.T) {
	t.Parallel()

	e2eTest := newE2ETest(t)
	assert.Assert(t, e2eTest != nil)

	r := test.NewKnRunResultCollector(t, e2eTest.it.KnTest())
	defer r.DumpIfFailed()

	err := e2eTest.it.KnPlugin().Install()
	assert.NilError(t, err)

	t.Log("kn-source_github create 'source-name' with 'sink-name'")
	e2eTest.ghSourceCreate(t, r, "source-name", "sink-name")

	t.Log("kn-source_github describe 'source-name'")
	e2eTest.ghSourceDescribe(t, r, "source-name")

	t.Log("kn-source_github update 'source-name' with 'new-sink-name'")
	e2eTest.ghSourceUpdate(t, r, "source-name", "new-sink-name")

	t.Log("kn-source_github delete 'source-name'")
	e2eTest.ghSourceDelete(t, r, "source-name", "sink-name")

	err = e2eTest.it.KnPlugin().Uninstall()
	assert.NilError(t, err)
}

// Private

func (et *e2eTest) ghSourceCreate(t *testing.T, r *test.KnRunResultCollector, sourceName, sinkName string) {
	out := et.it.KnPlugin().Run("create", sourceName, "--sink", sinkName)
	r.AssertNoError(out)
	assert.Check(t, util.ContainsAllIgnoreCase(out.Stdout, "create", sourceName, "namespace", et.it.KnTest().Namespace(), "sink", sinkName))
}

func (et *e2eTest) ghSourceDescribe(t *testing.T, r *test.KnRunResultCollector, sourceName string) {
	out := et.it.KnPlugin().Run("describe", sourceName)
	r.AssertNoError(out)
	assert.Check(t, util.ContainsAllIgnoreCase(out.Stdout, "describe", sourceName, "namespace", et.it.KnTest().Namespace()))
}

func (et *e2eTest) ghSourceUpdate(t *testing.T, r *test.KnRunResultCollector, sourceName, sinkName string) {
	out := et.it.KnPlugin().Run("update", sourceName, "--sink", sinkName)
	r.AssertNoError(out)
	assert.Check(t, util.ContainsAllIgnoreCase(out.Stdout, "update", sourceName, "namespace", et.it.KnTest().Namespace(), "sink", sinkName))
}

func (et *e2eTest) ghSourceDelete(t *testing.T, r *test.KnRunResultCollector, sourceName, sinkName string) {
	out := et.it.KnPlugin().Run("delete", sourceName)
	r.AssertNoError(out)
	assert.Check(t, util.ContainsAllIgnoreCase(out.Stdout, "delete", sourceName, "namespace", et.it.KnTest().Namespace()))
}
