/*Package cmd ...
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

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

	"github.com/spf13/cobra"
)

// patchCmd represents the patch command
var patchCmd = &cobra.Command{
	Use:              "patch",
	TraverseChildren: true,
	Short:            "Patch a resource in team namespace",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		identifyTeam()
	},
}

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:              "get",
	TraverseChildren: true,
	Short:            "Get resources in team namespace",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		identifyTeam()
	},
}

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:              "delete",
	TraverseChildren: true,
	Short:            "Delete a resource in team namespace",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		identifyTeam()
	},
}

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:              "create",
	TraverseChildren: true,
	Short:            "Create a resource in team namespace",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		identifyTeam()
	},
}

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("v%v\n", version)
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(deleteCmd)
	rootCmd.AddCommand(patchCmd)
	rootCmd.AddCommand(getCmd)
	rootCmd.AddCommand(versionCmd)
}
