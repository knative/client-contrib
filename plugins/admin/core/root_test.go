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
	"testing"

	"gotest.tools/assert"
)

func TestNewAdminCommand(t *testing.T) {

	t.Run("check subcommands", func(t *testing.T) {

		expectedSubCommands := []string{
			"version",
			"help",
			"domain",
			"registry",
		}

		cmd := NewAdminCommand()
		assert.Check(t, cmd.HasSubCommands())
		assert.Equal(t, len(cmd.Commands()), len(expectedSubCommands))

		for _, e := range expectedSubCommands {
			_, _, err := cmd.Find([]string{e})
			assert.NilError(t, err, "root command should have subcommand %q", e)

		}
	})

}
