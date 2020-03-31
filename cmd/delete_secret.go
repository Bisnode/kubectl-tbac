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

	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// deleteSecretCmd represents the deleteSecret command
var deleteSecretCmd = &cobra.Command{
	Use:   "secret [name]",
	Args:  cobra.ExactArgs(1),
	Short: "Delete a secret in your teams namespace.",
	Long: `
Delete a secret in your teams namespace. Your team is in the request if you are
logged in. If you belong to more than one team the command will ask you to provide
the --namespace flag.

Examples
# Delete a secret in your namespace with username and password.
kubectl tbac delete secret my-secret"

# Create a secret using namespace
kubectl tbac delete secret my-secret --namespace team-platform"
`,
	Run: func(cmd *cobra.Command, args []string) {
		deleteSecret(args)
	},
}

func deleteSecret(args []string) (err error) {
	secretName := &args[0]
	clientset, err := util.CreateClientSet()
	if err != nil {
		fmt.Printf("Failed to create client set: %v\n", err.Error())
		return err
	}

	// Delete the secret
	if err := clientset.
		CoreV1().
		Secrets(namespace).
		Delete(*secretName, &metav1.DeleteOptions{}); err != nil {
		fmt.Printf("Error deleting resource in namespace %v: %v\n", namespace, err.Error())
		return err
	}

	fmt.Printf("Deleted secret/%v in namespace %v\n", *secretName, namespace)
	return
}

func init() {
	deleteCmd.AddCommand(deleteSecretCmd)
}
