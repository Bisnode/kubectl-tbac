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
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// secretCmd represents the secret command
var createSecretCmd = &cobra.Command{
	Use:     "secret [name] [flags]",
	Aliases: secretAliases,
	Args:    cobra.ExactArgs(1),
	Short:   "Create a secret in your teams namespace",
	Long: `
Create a secret in your teams namespace. Your team is in the request if you are
logged in. If you belong to more than one team the command will ask you to provide
the --namespace flag.

Examples
# Create a secret in your namespace with username and password.
kubectl tbac create secret my-secret --data "USERNAME=foo" --data "PASSWORD=bar"

# Create a secret using namespace
kubectl tbac create secret my-secret --namespace team-platform -d "USER=foo" -d "PWD=bar"
`,

	Run: func(cmd *cobra.Command, args []string) {
		if err := createSecret(cmd, args); err != nil {
			fmt.Println(err)
		}
	},
}

func createSecret(cmd *cobra.Command, args []string) (err error) {
	secretName := &args[0]

	clientset, err := util.CreateClientSet()
	if err != nil {
		fmt.Printf("Failed to create client set: %v\n", err.Error())
		return err
	}

	secretsClient := clientset.CoreV1().Secrets(namespace)

	newSecret := &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      *secretName,
			Namespace: namespace,
			Labels: map[string]string{
				"app": *secretName,
			},
		},
		Data: util.AssembleInputData(data),
	}

	newSecret, err = secretsClient.Create(newSecret)
	if err != nil {
		fmt.Printf("Error creating resource: %v\n", err.Error())
		return err
	}

	fmt.Printf("Created secret/%v in namespace %v\n", newSecret.Name, namespace)
	return
}

func init() {
	createCmd.AddCommand(createSecretCmd)
	createSecretCmd.Flags().StringArrayVarP(&data, "data", "d", []string{}, "Data to add to secret")
}
