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
	"testing"

	"gotest.tools/assert"
)

func TestNewDomainCmd(t *testing.T) {
	cmd := NewDomainCmd(nil)
	assert.Check(t, cmd.HasSubCommands(), "cmd domain should have subcommands")
	assert.Equal(t, 3, len(cmd.Commands()), "domain command should have 3 subcommands")

	_, _, err := cmd.Find([]string{"set"})
	assert.NilError(t, err, "domain command should have set subcommand")

	_, _, err = cmd.Find([]string{"unset"})
	assert.NilError(t, err, "domain command should have unset subcommand")

	_, _, err = cmd.Find([]string{"help"})
	assert.NilError(t, err, "domain command should have help subcommand")
}
