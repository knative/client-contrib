/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package core

import (
	"fmt"
	"github.com/knative/client-contrib/plugins/kn-admin/cmd"
	"github.com/knative/client-contrib/plugins/kn-admin/cmd/domain"
	private_registry "github.com/knative/client-contrib/plugins/kn-admin/cmd/private-registry"
	"github.com/spf13/cobra"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands

func NewAdminCommand(params ...cmd.AdminParams) *cobra.Command {
	p := &cmd.AdminParams{}
	kubeConfig := os.Getenv("KUBECONFIG")
	if kubeConfig == "" {
		fmt.Println("cannot get cluster kube config, please export environment variable KUBECONFIG")
		os.Exit(1)
	}

	cfg, err := clientcmd.BuildConfigFromFlags("", kubeConfig)
	if err != nil {
		fmt.Println("failed to build config:", err)
		os.Exit(1)
	}

	p.ClientSet, err = kubernetes.NewForConfig(cfg)
	if err != nil {
		fmt.Println("failed to create client:", err)
		os.Exit(1)
	}

	rootCmd := &cobra.Command{
		Use:   "admin",
		Short: "A plugin of kn client to manage Knative",
		Long: `A plugin of kn client to manage Knative for administrators. 

For example:
kn admin domain set - to set Knative route domain
kn admin private-registry enable - to enable deployment from the private registry
`,
	}
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/kn/plugins/kn-admin.yaml)")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	//
	rootCmd.AddCommand(domain.NewDomainCmd(p))
	rootCmd.AddCommand(private_registry.NewPrivateRegistryCmd(p))
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

		// Search config in home directory with name ".kn-admin" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".kn-admin")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
