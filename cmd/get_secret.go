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
	"Bisnode/kubectl-tbac/util"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// getSecretCmd represents the getSecret command
var getSecretCmd = &cobra.Command{
	Use:     "secret [name]",
	Args:    cobra.RangeArgs(0, 1),
	Aliases: secretAliases,
	Short:   "Get a list of secrets or describe one.",
	Long:    `List secrets in team namespace or describe one.`,
	Run: func(cmd *cobra.Command, args []string) {
		getSecret(cmd, args)
	},
}

func getSecret(cmd *cobra.Command, args []string) (err error) {
	clientset, err := util.CreateClientSet()
	if err != nil {
		fmt.Printf("Failed to create client set: %v\n", err.Error())
		return err
	}

	listOpts := metav1.ListOptions{}
	if len(args) > 0 {
		listOpts = metav1.ListOptions{
			FieldSelector: fmt.Sprintf("metadata.name=%v", args[0]),
		}
	}

	secrets, err := clientset.
		CoreV1().
		Secrets(namespace).
		List(listOpts)

	if err != nil {
		fmt.Printf("Failed to list secrets in namespace %v: %v\n", namespace, err.Error())
		return err
	}

	// Describe if secret is specified. Otherwise list them all.
	if len(args) > 0 {
		fmt.Printf("Namespace:\t%v\n", namespace)
		fmt.Printf("Secret name:\t%v\n", secrets.Items[0].Name)
		fmt.Printf("Created:\t%v\n", secrets.Items[0].CreationTimestamp)
		if len(secrets.Items[0].Data) > 0 {
			fmt.Println(strings.Repeat("-", 46))
			for k, v := range secrets.Items[0].Data {
				fmt.Printf("%v:\n%v\n\n", k, string(v))
			}
		}
		return
	}
	for _, s := range secrets.Items {
		fmt.Printf(" * %v\n", s.Name)
	}
	return
}

func init() {
	getCmd.AddCommand(getSecretCmd)
}
