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
	"strings"

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
		// Handle shell completion: proxy __complete to kubectl
		if len(args) > 0 && (args[0] == "__complete" || args[0] == "__completeNoDesc") {
			return exec.KubeCtl(args)
		}

		// Parse --yes/-y flag from args before passing to kubectl
		args, confirmed := extractYesFlag(args)

		verb := findVerb(args)
		isSafe, err := safe.IsSafe(verb, args)
		if err != nil {
			return err
		}
		if !isSafe && !confirmed {
			if !prompt.Confirm(verb, args) {
				klog.Info("Not running command.")
				os.Exit(0)
			}
		}
		return exec.KubeCtl(args)
	},
}

// findVerb returns the first non-flag argument (the kubectl verb).
// It skips arguments starting with "-" and their values when they use
// the "--flag value" form (as opposed to "--flag=value").
func findVerb(args []string) string {
	flagsWithValues := map[string]bool{
		"--context": true, "-n": true, "--namespace": true,
		"--kubeconfig": true, "--cluster": true, "--user": true,
		"--as": true, "--as-group": true, "--certificate-authority": true,
		"--client-certificate": true, "--client-key": true, "--server": true,
		"--token": true, "-s": true,
	}
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if !strings.HasPrefix(arg, "-") {
			return arg
		}
		// --flag=value form: skip just this arg
		if strings.Contains(arg, "=") {
			continue
		}
		// --flag value form: skip this arg and the next
		if flagsWithValues[arg] && i+1 < len(args) {
			i++
		}
	}
	return ""
}

// extractYesFlag removes --yes or -y from the args (these are kubectl-safe flags,
// not kubectl flags) and returns the cleaned args and whether the flag was present.
func extractYesFlag(args []string) ([]string, bool) {
	found := false
	cleaned := make([]string, 0, len(args))
	for _, arg := range args {
		if arg == "--yes" || arg == "-y" {
			found = true
			continue
		}
		cleaned = append(cleaned, arg)
	}
	return cleaned, found
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
}
