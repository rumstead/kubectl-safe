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
	"os"

	"github.com/spf13/cobra"
	"k8s.io/klog/v2"

	"github.com/rumstead/kubectl-safe/pkg/cmd/safe"
	"github.com/rumstead/kubectl-safe/pkg/exec"
	"github.com/rumstead/kubectl-safe/pkg/prompt"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:                "safe",
	Args:               cobra.ArbitraryArgs,
	DisableFlagParsing: true,
	FParseErrWhitelist: cobra.FParseErrWhitelist{UnknownFlags: true},
	Short:              "A kubectl plugin to prevent shooting yourself in the foot with edit commands",
	Long:               "A kubectl plugin to prevent shooting yourself in the foot with edit commands",
	RunE: func(cmd *cobra.Command, args []string) error {
		verb := ""
		if len(args) > 0 {
			verb = args[0]
		}
		isSafe, err := safe.IsSafe(verb, args)
		if err != nil {
			return err
		}
		if !isSafe {
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
	return rootCmd.Execute()
}

func init() {
}
