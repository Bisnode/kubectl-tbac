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
	CreationTimestamp metav1.Time
	LastUpdated       string
	Service           string
	Container         string
	Data              map[string][]byte
}

// getSecretCmd represents the getSecret command
var getSecretCmd = &cobra.Command{
	Use:     "secret [name]",
	Args:    cobra.RangeArgs(0, 1),
	Aliases: secretAliases,
	Short:   "Get a list of secrets or describe one.",
	Long:    `List secrets in team namespace or describe one.`,
	Run: func(cmd *cobra.Command, args []string) {
		clientSet, err := util.CreateClientSet()
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
	data := make(map[string][]byte)
	for k, v := range secrets.Items[0].Data {
		data[k] = v
	}
	secretDesc = &SecretDescription{
		Namespace:         Namespace,
		Name:              secretName,
		LastUpdated:       secrets.Items[0].Annotations["tbac.bisnode.com/last-modified"],
		Service:           secrets.Items[0].Labels["app"],
		Container:         secrets.Items[0].Labels["tbac.bisnode.com/container"],
		CreationTimestamp: secrets.Items[0].CreationTimestamp,
		Data:              data,
	}
	return secretDesc, nil
}

func init() {
	getCmd.AddCommand(getSecretCmd)
}
