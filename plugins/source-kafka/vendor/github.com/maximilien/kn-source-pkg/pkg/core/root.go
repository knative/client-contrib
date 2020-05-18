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

package core

import (
	"github.com/maximilien/kn-source-pkg/pkg/types"

	"github.com/spf13/cobra"
)

func NewKnSourceCommand(knSourceFactory types.KnSourceFactory,
	commandFactory types.CommandFactory,
	flagsFactory types.FlagsFactory,
	runEFactory types.RunEFactory) *cobra.Command {
	knSourceParams := knSourceFactory.KnSourceParams()
	rootCmd := commandFactory.SourceCommand()

	// Disable docs header
	rootCmd.DisableAutoGenTag = true

	// Affects children as well
	rootCmd.SilenceUsage = true

	// Prevents Cobra from dealing with errors as we deal with them in main.go
	rootCmd.SilenceErrors = true

	if knSourceParams.Output != nil {
		rootCmd.SetOutput(knSourceParams.Output)
	}

	createCmd := commandFactory.CreateCommand()
	addCommonFlags(knSourceParams, createCmd)
	addCreateUpdateFlags(knSourceParams, createCmd)
	createCmd.Flags().AddFlagSet(flagsFactory.CreateFlags())
	createCmd.RunE = runEFactory.CreateRunE()
	rootCmd.AddCommand(createCmd)

	deleteCmd := commandFactory.DeleteCommand()
	addCommonFlags(knSourceParams, deleteCmd)
	deleteCmd.Flags().AddFlagSet(flagsFactory.DeleteFlags())
	deleteCmd.RunE = runEFactory.DeleteRunE()
	rootCmd.AddCommand(deleteCmd)

	updateCmd := commandFactory.UpdateCommand()
	if updateCmd != nil {
		addCommonFlags(knSourceParams, updateCmd)
		addCreateUpdateFlags(knSourceParams, updateCmd)
		updateCmd.Flags().AddFlagSet(flagsFactory.UpdateFlags())
		updateCmd.RunE = runEFactory.UpdateRunE()
		rootCmd.AddCommand(updateCmd)
	}

	describeCmd := commandFactory.DescribeCommand()
	addCommonFlags(knSourceParams, describeCmd)
	describeCmd.Flags().AddFlagSet(flagsFactory.DescribeFlags())
	describeCmd.RunE = runEFactory.DescribeRunE()
	rootCmd.AddCommand(describeCmd)

	// Initialize default `help` cmd early to prevent unknown command errors
	rootCmd.InitDefaultHelpCmd()

	return rootCmd
}

// Private

func addCommonFlags(knSourceParams *types.KnSourceParams, cmd *cobra.Command) {
	knSourceParams.AddCommonFlags(cmd)
}

func addCreateUpdateFlags(knSourceParams *types.KnSourceParams, cmd *cobra.Command) {
	knSourceParams.AddCreateUpdateFlags(cmd)
}
