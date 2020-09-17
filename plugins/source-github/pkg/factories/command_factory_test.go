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
)

func TestNewGHCommandFactory(t *testing.T) {
	knSourceFactory := NewGHSourceFactory()
	commandFactory := NewGHCommandFactory(knSourceFactory)

	assert.Assert(t, commandFactory != nil)
}

func TestSourceCommand(t *testing.T) {
	commandFactory := NewGHCommandFactory(NewGHSourceFactory())

	sourceCmd := commandFactory.SourceCommand()
	assert.Assert(t, sourceCmd != nil)

	assert.Equal(t, sourceCmd.Use, "github")
	assert.Equal(t, sourceCmd.Short, "Knative eventing GitHub source plugin")
	assert.Equal(t, sourceCmd.Long, "Manage your Knative GitHub eventing sources")
}

func TestCreateCommand(t *testing.T) {
	commandFactory := NewGHCommandFactory(NewGHSourceFactory())

	createCmd := commandFactory.CreateCommand()
	assert.Assert(t, createCmd != nil)

	assert.Equal(t, createCmd.Short, "create NAME")
	assert.Equal(t, createCmd.Long, "create a GitHub source")
	assert.Equal(t, createCmd.Example, `# Creates a new GitHub source with NAME using credentials
kn source github create NAME  --access-token $MY_ACCESS_TOKEN --secret-token $MY_SECRET_TOKEN

# Creates a new GitHub source with NAME with specified organization and repository using credentials
kn source github create NAME --org knative --repo client-contrib --access-token $MY_ACCESS_TOKEN --secret-token $MY_SECRET_TOKEN`)
}

func TestDeleteCommand(t *testing.T) {
	commandFactory := NewGHCommandFactory(NewGHSourceFactory())

	deleteCmd := commandFactory.DeleteCommand()
	assert.Assert(t, deleteCmd != nil)

	assert.Equal(t, deleteCmd.Short, "delete NAME")
	assert.Equal(t, deleteCmd.Long, "delete a GitHub source")
	assert.Equal(t, deleteCmd.Example, `# Deletes a GitHub source with NAME
kn source github delete NAME`)
}

func TestUpdateCommand(t *testing.T) {
	commandFactory := NewGHCommandFactory(NewGHSourceFactory())

	updateCmd := commandFactory.UpdateCommand()
	assert.Assert(t, updateCmd != nil)

	assert.Equal(t, updateCmd.Short, "update NAME")
	assert.Equal(t, updateCmd.Long, "update a GitHub source")
	assert.Equal(t, updateCmd.Example, `# Updates a GitHub source with NAME
kn source github update NAME`)
}

func TestDescribeCommand(t *testing.T) {
	commandFactory := NewGHCommandFactory(NewGHSourceFactory())

	describeCmd := commandFactory.DescribeCommand()
	assert.Assert(t, describeCmd != nil)

	assert.Equal(t, describeCmd.Short, "describe NAME")
	assert.Equal(t, describeCmd.Long, "update a GitHub source")
	assert.Equal(t, describeCmd.Example, `# Describes a GitHub source with NAME
kn source github describe NAME`)
}
