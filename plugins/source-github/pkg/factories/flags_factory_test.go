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

package factories

import (
	"testing"

	"gotest.tools/assert"

	"github.com/spf13/pflag"
)

func TestNewGHFlagsFactory(t *testing.T) {
	flagsFactory := NewGHFlagsFactory(NewGHSourceFactory())
	assert.Assert(t, flagsFactory != nil)
}

func TestCreateFlags(t *testing.T) {
	flagsFactory := NewGHFlagsFactory(NewGHSourceFactory())

	createFlags := flagsFactory.CreateFlags()
	assert.Assert(t, createFlags != nil)

	testCreateUpdateFlags(t, createFlags)
}

func TestDeleteFlags(t *testing.T) {
	flagsFactory := NewGHFlagsFactory(NewGHSourceFactory())

	deleteFlags := flagsFactory.DescribeFlags()
	assert.Assert(t, deleteFlags != nil)
}

func TestUpdateFlags(t *testing.T) {
	flagsFactory := NewGHFlagsFactory(NewGHSourceFactory())

	updateFlags := flagsFactory.CreateFlags()
	assert.Assert(t, updateFlags != nil)

	testCreateUpdateFlags(t, updateFlags)
}

func TestDescribeFlags(t *testing.T) {
	flagsFactory := NewGHFlagsFactory(NewGHSourceFactory())

	describeFlags := flagsFactory.DescribeFlags()
	assert.Assert(t, describeFlags != nil)
}

// Private

func testCreateUpdateFlags(t *testing.T, flagSet *pflag.FlagSet) {
	orgFlag, err := flagSet.GetString("org")
	assert.NilError(t, err)
	assert.Assert(t, orgFlag == "")

	repoFlag, err := flagSet.GetString("repo")
	assert.NilError(t, err)
	assert.Assert(t, repoFlag == "")

	apiURLFlag, err := flagSet.GetString("api-url")
	assert.NilError(t, err)
	assert.Assert(t, apiURLFlag == "https://api.github.com")

	secretTokenFlag, err := flagSet.GetString("secret-token")
	assert.NilError(t, err)
	assert.Assert(t, secretTokenFlag == "")

	accessTokenFlag, err := flagSet.GetString("access-token")
	assert.NilError(t, err)
	assert.Assert(t, accessTokenFlag == "")
}
