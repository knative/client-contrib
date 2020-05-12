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
	"github.com/maximilien/kn-source-pkg/pkg/types"
	"knative.dev/client/pkg/kn/commands"

	"github.com/spf13/pflag"
)

type DefautFlagsFactory struct {
	knSourceFactory types.KnSourceFactory
}

func NewDefaultFlagsFactory(knSourceFactory types.KnSourceFactory) types.FlagsFactory {
	return &DefautFlagsFactory{
		knSourceFactory: knSourceFactory,
	}
}

func (f *DefautFlagsFactory) CreateFlags() *pflag.FlagSet {
	flagSet := pflag.NewFlagSet("create", pflag.ExitOnError)
	f.addNamespaceFlag(flagSet)
	return flagSet
}

func (f *DefautFlagsFactory) DeleteFlags() *pflag.FlagSet {
	flagSet := pflag.NewFlagSet("delete", pflag.ExitOnError)
	f.addNamespaceFlag(flagSet)
	return flagSet
}

func (f *DefautFlagsFactory) UpdateFlags() *pflag.FlagSet {
	flagSet := pflag.NewFlagSet("update", pflag.ExitOnError)
	f.addNamespaceFlag(flagSet)
	return flagSet
}

func (f *DefautFlagsFactory) DescribeFlags() *pflag.FlagSet {
	flagSet := pflag.NewFlagSet("describe", pflag.ExitOnError)
	f.addNamespaceFlag(flagSet)
	return flagSet
}

// Private

func (f *DefautFlagsFactory) addNamespaceFlag(flagSet *pflag.FlagSet) {
	commands.AddNamespaceFlags(flagSet, false)
}
