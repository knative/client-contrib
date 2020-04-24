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

	"github.com/spf13/cobra"
	"knative.dev/client-contrib/plugins/migration/pkg/command"
	"knative.dev/client-contrib/plugins/migration/pkg/command/list"
	"knative.dev/client-contrib/plugins/migration/pkg/command/migrate"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string

// migrationCmd represents the base command when called without any subcommands
func NewMigrationCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "kn\u00A0migration",
		Short: "A plugin of kn client to migrate Knative services",
		Long: `kn admin: a plugin of kn client to migrate Knative services.
For example:
kn migration list
kn migration migrate --namespace default --destination-namespace default
`,
	}
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/kn/plugins/admin.yaml)")
	rootCmd.AddCommand(list.NewListCommand())
	rootCmd.AddCommand(migrate.NewMigrateCommand())
	rootCmd.AddCommand(command.NewVersionCommand())
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

		// Search config in home directory with name ".migration" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".migration")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
