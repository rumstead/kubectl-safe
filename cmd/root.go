// Package cmd
/*
Copyright © 2022 Ryan Umstead rjumstead@gmail.com

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
package cmd

import (
	"github.com/rumstead/kubectl-safe/pkg/cmd/safe"
	"github.com/rumstead/kubectl-safe/pkg/exec"
	"github.com/rumstead/kubectl-safe/pkg/prompt"
	"github.com/spf13/cobra"
	"k8s.io/klog/v2"
	"os"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:                "safe",
	Args:               cobra.ArbitraryArgs,
	DisableFlagParsing: true,
	FParseErrWhitelist: cobra.FParseErrWhitelist{UnknownFlags: true},
	Short:              "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		verb := ""
		if len(args) > 0 {
			verb = args[0]
		}
		safe, err := safe.IsVerbSafe(verb)
		if err != nil {
			return err
		}
		if !safe {
			if !prompt.Confirm(verb) {
				klog.Info("Not running command.")
				os.Exit(0)
			}
		}
		return exec.KubeCtl(args)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	err := rootCmd.Execute()
	if err != nil {
		return err
	}
	return nil
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.kubectl-safe.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
