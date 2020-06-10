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

package main

import (
	"fmt"
	"os"

	"github.com/maximilien/kn-source-github/pkg/factories"
	
	"github.com/maximilien/kn-source-pkg/pkg/core"
	"github.com/maximilien/kn-source-pkg/pkg/util"
)

func main() {
	ghSourceFactory := factories.NewGHSourceFactory()

	ghCommandFactory := factories.NewGHCommandFactory(ghSourceFactory)
	ghFlagsFactory := factories.NewGHFlagsFactory(ghSourceFactory)
	ghRunEFactory := factories.NewGHRunEFactory(ghSourceFactory)

	rootCmd := core.NewKnSourceCommand(ghSourceFactory, ghCommandFactory, ghFlagsFactory, ghRunEFactory)
	err := util.ReadmeGenerator(rootCmd)
	if err != nil {
		fmt.Fprintln(os.Stdout, err)
		os.Exit(1)
	}
}