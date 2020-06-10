// Copyright © 2018 The Knative Authors
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
	"github.com/maximilien/kn-source-github/pkg/types"

	sourcefactories "github.com/maximilien/kn-source-pkg/pkg/factories"
	sourcetypes "github.com/maximilien/kn-source-pkg/pkg/types"

	"github.com/spf13/cobra"
)

type ghCommandFactory struct {
	ghSourceFactory       types.GHSourceFactory
	defaultCommandFactory sourcetypes.CommandFactory
}

func NewGHCommandFactory(ghSourceFactory types.GHSourceFactory) types.GHCommandFactory {
	return &ghCommandFactory{
		ghSourceFactory:       ghSourceFactory,
		defaultCommandFactory: sourcefactories.NewDefaultCommandFactory(ghSourceFactory),
	}
}

func (f *ghCommandFactory) SourceCommand() *cobra.Command {
	sourceCmd := f.defaultCommandFactory.SourceCommand()
	sourceCmd.Use = "github"
	sourceCmd.Short = "Knative eventing GitHub source plugin"
	sourceCmd.Long = "Manage your Knative GitHub eventing sources"
	return sourceCmd
}

func (f *ghCommandFactory) CreateCommand() *cobra.Command {
	createCmd := f.defaultCommandFactory.CreateCommand()
	createCmd.Short = "create NAME"
	createCmd.Long = "create a GitHub source"
	createCmd.Example = `# Creates a new GitHub source with NAME using credentials
kn source github create NAME  --access-token $MY_ACCESS_TOKEN --secret-token $MY_SECRET_TOKEN

# Creates a new GitHub source with NAME with specified organization and repository using credentials
kn source github create NAME --org knative --repo client-contrib --access-token $MY_ACCESS_TOKEN --secret-token $MY_SECRET_TOKEN`
	return createCmd
}

func (f *ghCommandFactory) DeleteCommand() *cobra.Command {
	deleteCmd := f.defaultCommandFactory.DeleteCommand()
	deleteCmd.Short = "delete NAME"
	deleteCmd.Long = "delete a GitHub source"
	deleteCmd.Example = `# Deletes a GitHub source with NAME
kn source github delete NAME`
	return deleteCmd
}

func (f *ghCommandFactory) UpdateCommand() *cobra.Command {
	updateCmd := f.defaultCommandFactory.UpdateCommand()
	updateCmd.Short = "update NAME"
	updateCmd.Long = "update a GitHub source"
	updateCmd.Example = `# Updates a GitHub source with NAME
kn source github update NAME`
	return updateCmd
}

func (f *ghCommandFactory) DescribeCommand() *cobra.Command {
	describeCmd := f.defaultCommandFactory.DescribeCommand()
	describeCmd.Short = "describe NAME"
	describeCmd.Long = "update a GitHub source"
	describeCmd.Example = `# Describes a GitHub source with NAME
kn source github describe NAME`
	return describeCmd
}
