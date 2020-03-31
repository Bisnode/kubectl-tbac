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
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

var removeData []string

// patchSecretCmd represents the patchSecret command
var patchSecretCmd = &cobra.Command{
	Use:     "secret [name] [--data key=value|--remove-data key]",
	Aliases: secretAliases,
	Args:    cobra.ExactArgs(1),
	Short:   "Patch a secret in your teams namespace",
	Long: `
Patches a secret in your teams namespaced. Your team is in the request if you are
logged in. If you belong to more than one team the command will ask you to provide
the --namespace flag.

Examples
# Patch a secret in your namespace with username and password.
kubectl tbac patch secret my-secret --data "USERNAME=foo" --data "PASSWORD=bar"

# Patch a secret using namespace
kubectl tbac patch secret my-secret --namespace team-platform -d "USER=foo" -d "PWD=bar"

# Remove secret key USERNAME and PASSWORD from secret
kubectl tbac patch secret my-secret --remove-data USERNAME --remove-data PASSWORD
`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := patchSecret(cmd, args); err != nil {
			fmt.Println(err)
		}
	},
}

func patchSecret(cmd *cobra.Command, args []string) (err error) {
	secretName := &args[0]

	// Create secrets client
	clientset, err := util.CreateClientSet()
	if err != nil {
		return err
	}
	secretsClient := clientset.CoreV1().Secrets(namespace)

	// Get current secret
	originalSecret, err := secretsClient.Get(*secretName, metav1.GetOptions{})
	if err != nil {
		return err
	}
	patchSecret := originalSecret

	// Resource version cannot be defined in request to Kubernetes.
	originalSecret.ResourceVersion = ""
	patchSecret.ResourceVersion = ""

	// If data to remove, remove from the patch version
	for _, d := range removeData {
		if originalSecret.Data[d] != nil {
			delete(patchSecret.Data, d)
		}
	}

	/*
		When removing keys in the data of a secret using the --remove-data flag,
		the secret is first removed from Kubernetes and then recreated
		without the unwanted keys.
	*/
	if len(removeData) != 0 {
		if err := secretsClient.Delete(*secretName, &metav1.DeleteOptions{}); err != nil {
			return err
		}
		patchSecret, err = secretsClient.Create(patchSecret)
		/*If delete was successful but recreate not:
		- Try to roll back to original secret and abort.
		- If rollback is unsuccessful then dump data to terminal and abort.*/
		if err != nil {
			fmt.Printf("Secret recreation failed: %v\n", err)
			fmt.Println("Attempt to roll back to original secret")
			_, err := secretsClient.Create(originalSecret)
			if err != nil {
				fmt.Printf("Rollback failed: %v\n", err)
				fmt.Printf(`The original secret %v was removed and could not be recreated
					Recreation need to be handled manually. It contained data:`, *secretName)
				if len(originalSecret.Data) == 0 {
					fmt.Printf("No data\n")
				}
				for k, v := range originalSecret.Data {
					fmt.Printf("%v:\n%v\n\n", k, string(v))
				}
			}
			return fmt.Errorf("Errors during recreation - Can't continue")
		}
	}

	for k, v := range util.AssembleInputData(data) {
		patchSecret.Data[k] = v
	}

	patch, err := json.Marshal(patchSecret)
	if err != nil {
		return err
	}

	_, err = secretsClient.Patch(*secretName, types.StrategicMergePatchType, patch)
	if err != nil {
		return err
	}
	fmt.Printf("secret/%v modified\n", *secretName)
	return
}

func init() {
	patchCmd.AddCommand(patchSecretCmd)
	patchSecretCmd.Flags().StringArrayVarP(&data, "data", "d", []string{}, "Data to add or update in secret")
	patchSecretCmd.Flags().StringArrayVarP(&removeData, "remove-data", "r", []string{}, "Remove data key from secret")
}
