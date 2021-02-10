package cmd

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

// GenerateSecrets is a set of secret definitions.
var GenerateSecrets = []v1.Secret{
	v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-credentials",
			Namespace: "default",
			Labels: map[string]string{
				"app":                        "my-credentials",
				"tbac.bisnode.com/container": "default",
			},
			Annotations: map[string]string{
				"tbac.bisnode.com/last-modified": fmt.Sprintf("%v", metav1.Now().Rfc3339Copy()),
			},
		},
		Data: map[string][]byte{
			"USERNAME": []byte("foo"),
			"PASSWORD": []byte("bar"),
			"KEY":      []byte("extra-key"),
		},
	},
	v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-api-key",
			Namespace: "default",
			Labels: map[string]string{
				"app":                        "my-api-key",
				"tbac.bisnode.com/container": "default",
			},
			Annotations: map[string]string{
				"tbac.bisnode.com/last-modified": fmt.Sprintf("%v", metav1.Now().Rfc3339Copy()),
			},
		},
		Data: map[string][]byte{
			"URL": []byte("github.com"),
			"KEY": []byte("key"),
		},
	},
}

func TestGetSecretsList(t *testing.T) {
	// create the 'fake' clientSet where clientSet.Interface = &Clientset{}, setting all the 'fake' methods
	// as seen in https://github.com/kubernetes/client-go/blob/master/kubernetes/fake/clientSet_generated.go
	clientSet := fake.NewSimpleClientset()
	Namespace = "default"

	for _, s := range GenerateSecrets {
		_, err := clientSet.CoreV1().Secrets(Namespace).Create(&s)
		if err != nil {
			assert.Equal(t, nil, err)
		}
	}

	secretList, err := GetSecretList(clientSet)
	if err != nil {
		assert.Equal(t, nil, err)
	}
	assert.Equal(t, len(GenerateSecrets), len(secretList))
	assert.Contains(t, secretList, "my-credentials")
	assert.Contains(t, secretList, "my-api-key")
}

func TestDescribeOneSecret(t *testing.T) {
	clientSet := fake.NewSimpleClientset()
	Namespace = "default"

	for _, s := range GenerateSecrets {
		_, err := clientSet.CoreV1().Secrets(Namespace).Create(&s)
		if err != nil {
			assert.Equal(t, nil, err)
		}
	}

	secretDescription, err := GetSecretDescription(clientSet, "my-credentials")
	if err != nil {
		assert.Equal(t, nil, err)
	}
	assert.Equal(t, "my-credentials", secretDescription.Name)
	assert.Equal(t, []byte("foo"), secretDescription.Data["USERNAME"])
	assert.Equal(t, []byte("bar"), secretDescription.Data["PASSWORD"])
}

func TestDeleteSecret(t *testing.T) {
	clientSet := fake.NewSimpleClientset()
	Namespace = "default"

	for _, s := range GenerateSecrets {
		_, err := clientSet.CoreV1().Secrets(Namespace).Create(&s)
		if err != nil {
			assert.Equal(t, nil, err)
		}
	}

	err := DeleteSecret(clientSet, "my-credentials")
	if err != nil {
		assert.Equal(t, nil, err)
	}
	secretList, err := clientSet.CoreV1().Secrets(Namespace).List(metav1.ListOptions{})
	if err != nil {
		assert.Equal(t, nil, err)
	}
	secretNames := []string{}
	for _, secret := range secretList.Items {
		secretNames = append(secretNames, secret.Name)
	}
	assert.Contains(t, secretNames, "my-api-key")
	assert.NotContains(t, secretNames, "my-credentials")
}

func TestPatchSecret(t *testing.T) {
	clientSet := fake.NewSimpleClientset()
	Namespace = "default"

	for _, s := range GenerateSecrets {
		_, err := clientSet.CoreV1().Secrets(Namespace).Create(&s)
		if err != nil {
			assert.Equal(t, nil, err)
		}
	}

	secretName := "my-credentials"
	removeData := []string{"USERNAME"}
	updateData := []string{
		"PASSWORD=snowmobile2020",
		"URL=my-api.com",
	}

	err := PatchSecret(clientSet, &secretName, &removeData, &updateData)
	if err != nil {
		assert.Equal(t, nil, err)
	}

	updatedSecret, err := clientSet.CoreV1().Secrets(Namespace).Get("my-credentials", metav1.GetOptions{})
	if err != nil {
		assert.Equal(t, nil, err)
	}

	// Should have updated
	assert.Equal(t, []byte("snowmobile2020"), updatedSecret.Data["PASSWORD"])
	assert.Equal(t, []byte("my-api.com"), updatedSecret.Data["URL"])
	// Should be intact
	assert.Equal(t, []byte("extra-key"), updatedSecret.Data["KEY"])
	// Should be gone
	assert.NotContains(t, updatedSecret.Data, "USERNAME")
}

func TestCreateSecret(t *testing.T) {
	clientSet := fake.NewSimpleClientset()
	Namespace = "default"

	secretName := "new-app-secret"
	container := "default"
	data := []string{
		"USERNAME=foo",
		"PASSWORD=bar",
	}
	err := CreateSecret(clientSet, &secretName, &container, data)
	if err != nil {
		assert.Equal(t, nil, err)
	}
	createdSecret, err := clientSet.CoreV1().Secrets(Namespace).Get("new-app-secret-default", metav1.GetOptions{})
	if err != nil {
		assert.Equal(t, nil, err)
	}
	assert.Equal(t, []byte("foo"), createdSecret.Data["USERNAME"])
	assert.Equal(t, []byte("bar"), createdSecret.Data["PASSWORD"])
}
