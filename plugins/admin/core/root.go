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
	"fmt"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"knative.dev/client-contrib/plugins/admin/pkg"
	"knative.dev/client-contrib/plugins/admin/pkg/command"
	"knative.dev/client-contrib/plugins/admin/pkg/command/autoscaling"
	"knative.dev/client-contrib/plugins/admin/pkg/command/domain"
	"knative.dev/client-contrib/plugins/admin/pkg/command/profiling"
	private_registry "knative.dev/client-contrib/plugins/admin/pkg/command/registry"
)

var cfgFile string

// NewAdminCommand represents the base command when called without any subcommands
func NewAdminCommand(params ...pkg.AdminParams) *cobra.Command {
	p := &pkg.AdminParams{}
	p.Initialize()

	rootCmd := &cobra.Command{
		Use:   "kn\u00A0admin",
		Short: "A plugin of kn client to manage Knative",
		Long:  `kn admin: a plugin of kn client to manage Knative for administrators`,

		// disable printing usage when error occurs
		SilenceUsage: true,
	}
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/kn/plugins/admin.yaml)")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.SetOut(os.Stdout)
	rootCmd.AddCommand(domain.NewDomainCmd(p))
	rootCmd.AddCommand(private_registry.NewPrivateRegistryCmd(p))
	rootCmd.AddCommand(autoscaling.NewAutoscalingCmd(p))
	rootCmd.AddCommand(profiling.NewProfilingCommand(p))
	rootCmd.AddCommand(command.NewVersionCommand())

	// Add default help page if there's unknown command
	rootCmd.InitDefaultHelpCmd()
	return rootCmd
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".admin" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".admin")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
