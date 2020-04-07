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

package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"knative.dev/client-contrib/plugins/hello/pkg/command"
)

/**
 * Sample main which just prints out a friendly message
 */

var rootCmd = &cobra.Command{
	Use:   "kn-hello",
	Short: "Sample kn plugin printing out a nice message",
	Long:  `Longer description of this fantastic plugin that can go over several lines.`,
}

func init() {
	rootCmd.AddCommand(command.NewPrintCommand())
	rootCmd.AddCommand(command.NewVersionCommand())
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		os.Exit(1)
	}
}
