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

package source

import (
	"github.com/maximilien/kn-source-pkg/pkg/types"
	"github.com/spf13/cobra"
)

// NewSourceCommand as the root group command
func NewSourceCommand(knSourceParams *types.KnSourceParams) *cobra.Command {
	createCmd := &cobra.Command{
		Use:   "source",
		Short: "Knative eventing {{.Name}} source plugin",
		Long:  "Manage your Knative {{.Name}} eventing sources",
	}
	return createCmd
}
