/*Package cmd ...
Copyright © 2020 Daniel Olsson

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
	"fmt"
	"log"
	"os"

	"github.com/Bisnode/kubectl-tbac/util"
	"github.com/spf13/cobra"
)

// Namespace in kubernetes.
var Namespace string

// Context name
var Context string

var (
	cfgFile string
	verbose bool
	lab     bool
	sandbox bool
	teams   []string
	data    []string
)

var secretAliases = []string{
	"sec",
	"secr",
	"secre",
	"secrets",
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "kubectl-tbac",
	Short: "Simplify managing resources based on tbac.",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().BoolVarP(&lab, "lab", "", false, "Run lab to simulate team membership.")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().BoolVarP(&sandbox, "sandbox", "s", false, "Set if you want to work in a sandbox Namespace.")
	rootCmd.PersistentFlags().StringVarP(&Namespace, "namespace", "n", "", "Namespace to create secret in. Usually only needed when member of more than one team.")
	rootCmd.PersistentFlags().StringVarP(&Context, "context", "", "", "Set context name.")

	// Hide flags
	rootCmd.PersistentFlags().MarkHidden("lab")
}

/* ---------------------------------------------
Custom functions to be used for all commands
Still not that generic that they fit in utils.
*/

// identifyTeam sets namespace based on team in access token.
// If sandbox is set, then appending namespace with "-sandbox"
func identifyTeam() {
	if lab {
		teams = []string{"team-platform"}
	} else {
		matchPrefix := "sec-tbac-team-"
		trimPrefix := "sec-tbac-"
		teams = util.WhoAmI(&matchPrefix, &trimPrefix, &Context)
	}

	if len(teams) == 1 {
		Namespace = teams[0]
	}

	if Namespace == "" && len(teams) > 1 {
		fmt.Println("You are member of multiple teams. Please use --namespace [team-name] to specify which namespace you want to work in.")
		for _, team := range teams {
			fmt.Println("-", team)
		}
		os.Exit(1)
	}
	if sandbox {
		Namespace = Namespace + "-sandbox"
	}
}
