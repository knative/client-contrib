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
//
//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

package types

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"k8s.io/client-go/rest"
)

// RunE abstracts the Cobra RunE interface into a usable type
type RunE = func(cmd *cobra.Command, args []string) error

// KnSourceClient is the base interface for all kn-source-extension
//counterfeiter:generate . KnSourceClient
type KnSourceClient interface {
	KnSourceParams() *KnSourceParams
	Namespace() string
	RestConfig() *rest.Config
}

// KnSourceFactory is the base factory interface for all kn-source-extension factories
//counterfeiter:generate . KnSourceFactory
type KnSourceFactory interface {
	KnSourceParams() *KnSourceParams

	CreateKnSourceParams() *KnSourceParams
	CreateKnSourceClient(restConfig *rest.Config, namespace string) KnSourceClient
}

// CommandFactory is the factory for cobra.Command objects
//counterfeiter:generate . CommandFactory
type CommandFactory interface {
	SourceCommand() *cobra.Command

	CreateCommand() *cobra.Command
	DeleteCommand() *cobra.Command
	UpdateCommand() *cobra.Command
	DescribeCommand() *cobra.Command

	KnSourceFactory() KnSourceFactory
}

// FlagsFactory is the factory for pflag.FlagSet objects
//counterfeiter:generate . FlagsFactory
type FlagsFactory interface {
	CreateFlags() *pflag.FlagSet
	DeleteFlags() *pflag.FlagSet
	UpdateFlags() *pflag.FlagSet
	DescribeFlags() *pflag.FlagSet

	KnSourceFactory() KnSourceFactory
}

// RunEFactory is the factory for RunE objects
//counterfeiter:generate . RunEFactory
type RunEFactory interface {
	CreateRunE() RunE
	DeleteRunE() RunE
	UpdateRunE() RunE
	DescribeRunE() RunE

	KnSourceFactory() KnSourceFactory
	KnSourceClient(restConfig *rest.Config, namespace string) KnSourceClient
}
