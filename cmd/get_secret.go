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
	"strings"

	"github.com/Bisnode/kubectl-tbac/util"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// SecretDescription holds data needed to describe a secret.
type SecretDescription struct {
	Namespace         string
	Name              string
	CreationTimestamp string
	LastUpdated       string
	Service           string
	Container         string
	Data              map[string][]byte
}

var export bool

// getSecretCmd represents the getSecret command
var getSecretCmd = &cobra.Command{
	Use:     "secret [name]",
	Args:    cobra.RangeArgs(0, 1),
	Aliases: secretAliases,
	Short:   "Get a list of secrets or describe one.",
	Long:    `List secrets in team namespace or describe one.`,
	Run: func(cmd *cobra.Command, args []string) {
		clientSet, err := util.CreateClientSet(&Context)
		if err != nil {
			fmt.Printf("Failed to create clientSet: %v\n", err)
			os.Exit(1)
		}
		if len(args) == 1 {
			secretDesc, err := GetSecretDescription(clientSet, args[0])
			if err != nil {
				fmt.Printf("Failed to get secret %v: %v", args[0], err)
				os.Exit(1)
			}
			if export {
				secretDesc.ExportSecret()
				os.Exit(0)
			}
			secretDesc.PrettyPrintSecretDesc()
			os.Exit(0)
		}
		secretList, err := GetSecretList(clientSet)
		if err != nil {
			fmt.Printf("Failed to get secrets: %v", err)
			os.Exit(1)
		}
		for _, s := range secretList {
			fmt.Printf(" * %v\n", s)
		}
	},
}

// PrettyPrintSecretDesc pretty prints a secret as a table view
func (s *SecretDescription) PrettyPrintSecretDesc() {
	fmt.Printf("Secret name:%v%v\n", strings.Repeat(" ", 25-len("Secret Name:")), s.Name)
	fmt.Printf("Service (app label):%v%v\n", strings.Repeat(" ", 25-len("Service (app label):")), s.Service)
	fmt.Printf("Container:%v%v\n", strings.Repeat(" ", 25-len("Container:")), s.Container)
	fmt.Printf("Namespace:%v%v\n", strings.Repeat(" ", 25-len("Namespace:")), s.Namespace)
	fmt.Printf("Created:%v%v\n", strings.Repeat(" ", 25-len("Created:")), s.CreationTimestamp)
	fmt.Printf("Last updated:%v%v\n\n", strings.Repeat(" ", 25-len("Last updated:")), s.LastUpdated)
	if len(s.Data) > 0 {
		fmt.Println(strings.Repeat("-", 25), "DATA", strings.Repeat("-", 25))
		for k, v := range s.Data {
			fmt.Printf("%v=%v\n", k, string(v))
		}
	}
}

// ExportSecret prints out secret in exported format.
func (s *SecretDescription) ExportSecret() {
	// Trim away container from name if found..

	name := strings.TrimSuffix(s.Name, "-"+s.Container)

	var d []string
	for k, v := range s.Data {
		d = append(d, fmt.Sprintf("--data '%v=%v'", string(k), string(v)))
	}
	data := strings.Join(d, " ")
	out := fmt.Sprintf("kubectl tbac create secret %v --namespace %v %v", name, s.Namespace, data)
	fmt.Printf("%v\n\n", out)
}

// GetSecretList returns a list of secrets in the namespace
func GetSecretList(clientSet kubernetes.Interface) (secrets []string, err error) {
	secretList, err := clientSet.
		CoreV1().
		Secrets(Namespace).
		List(metav1.ListOptions{
			FieldSelector: fmt.Sprintf("type=Opaque"),
		})

	if err != nil {
		fmt.Printf("Failed to list secrets in namespace %v: %v\n", Namespace, err.Error())
		return nil, err
	}
	for _, s := range secretList.Items {
		secrets = append(secrets, s.Name)
	}

	if len(secretList.Items) == 0 {
		fmt.Println("No resources found.")
	}
	return secrets, nil
}

// GetSecretDescription takes a secret name as input and return it in a SecretDescription.
func GetSecretDescription(clientSet kubernetes.Interface, secretName string) (secretDesc *SecretDescription, err error) {
	listOpts := metav1.ListOptions{
		FieldSelector: fmt.Sprintf("metadata.name=%v", secretName),
	}
	secrets, err := clientSet.
		CoreV1().
		Secrets(Namespace).
		List(listOpts)

	if err != nil {
		fmt.Printf("Failed to list secrets in namespace %v: %v\n", Namespace, err.Error())
		return nil, err
	}

	// Due to an issue with the ListOptions filters in
	// kubernetes fake-client that is used for testing
	// we cannot check for exactly one result
	if len(secrets.Items) < 1 {
		err := fmt.Errorf("Secret not found: %v/%v", Namespace, secretName)
		return nil, err
	}

	data := make(map[string][]byte)
	for k, v := range secrets.Items[0].Data {
		data[k] = v
	}
	secretDesc = &SecretDescription{
		Namespace:         Namespace,
		Name:              secretName,
		LastUpdated:       secrets.Items[0].Annotations["tbac.bisnode.com/last-modified"],
		CreationTimestamp: secrets.Items[0].Annotations["tbac.bisnode.com/time-created"],
		Service:           secrets.Items[0].Labels["app"],
		Container:         secrets.Items[0].Labels["tbac.bisnode.com/container"],
		Data:              data,
	}
	return secretDesc, nil
}

func init() {
	getCmd.AddCommand(getSecretCmd)
	getSecretCmd.PersistentFlags().BoolVarP(&export, "export", "", false, "Export as a `kubectl create secret` command")
}
