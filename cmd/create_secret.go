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
	"os"

	"github.com/Bisnode/kubectl-tbac/util"
	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var container string
var app string

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

# Create a secret for a sidecar named opa
kubectl tbac create secret my-secret --container opa -d "USER=foo" -d "PWD=bar"
`,

	Run: func(cmd *cobra.Command, args []string) {
		clientSet, err := util.CreateClientSet(&Context)
		if err != nil {
			fmt.Printf("Failed to create clientSet: %v\n", err)
			os.Exit(1)
		}
		if err := CreateSecret(clientSet, &args[0], &container, data); err != nil {
			fmt.Println(err)
		}
	},
}

// CreateSecret creates a secret in teams namespace
func CreateSecret(clientSet kubernetes.Interface, secretName, container *string, data []string) (err error) {
	secretsClient := clientSet.CoreV1().Secrets(Namespace)
	appLabel := *secretName

	if app != "" {
		appLabel = app
	}

	newSecret := &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      *secretName + "-" + *container,
			Namespace: Namespace,
			Labels: map[string]string{
				"app":                        appLabel,
				"tbac.bisnode.com/container": *container,
				"tbac.bisnode.com/sandbox":   fmt.Sprintf("%v", sandbox),
			},
			Annotations: map[string]string{
				"tbac.bisnode.com/last-modified": fmt.Sprintf("%v", metav1.Now().Rfc3339Copy()),
				"tbac.bisnode.com/time-created":  fmt.Sprintf("%v", metav1.Now().Rfc3339Copy()),
			},
		},
		Data: util.AssembleInputData(data),
	}

	newSecret, err = secretsClient.Create(newSecret)
	if err != nil {
		fmt.Printf("Error creating resource: %v\n", err.Error())
		return err
	}

	fmt.Printf("Created secret/%v in namespace %v\n", newSecret.Name, Namespace)
	return
}

func init() {
	createCmd.AddCommand(createSecretCmd)
	createSecretCmd.Flags().StringArrayVarP(&data, "data", "d", []string{}, "Data to add to secret")
	createSecretCmd.Flags().StringVarP(&container, "container", "c", "default", "Which container to create secret for. Only set this if you want to create a secret for a sidecar. (Default: \"default\"")
	createSecretCmd.Flags().StringVarP(&app, "app", "a", "", "Set the app label different than the secret name. Note that the app label must match the app label on the service that should use this secret.")
}
