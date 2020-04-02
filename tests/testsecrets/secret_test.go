package testsecrets

import (
	"testing"

	"github.com/mdanielolsson/kubectl-tbac/cmd"
	"github.com/mdanielolsson/kubectl-tbac/tests/testdata"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestGetSecretsList(t *testing.T) {
	// create the 'fake' clientSet where clientSet.Interface = &Clientset{}, setting all the 'fake' methods
	// as seen in https://github.com/kubernetes/client-go/blob/master/kubernetes/fake/clientSet_generated.go
	clientSet := fake.NewSimpleClientset()
	cmd.Namespace = "default"

	for _, s := range testdata.GenerateSecrets {
		_, err := clientSet.CoreV1().Secrets(cmd.Namespace).Create(&s)
		if err != nil {
			assert.Equal(t, nil, err)
		}
	}

	secretList, err := cmd.GetSecretList(clientSet)
	if err != nil {
		assert.Equal(t, nil, err)
	}
	assert.Equal(t, len(testdata.GenerateSecrets), len(secretList))
	assert.Contains(t, secretList, "my-credentials")
	assert.Contains(t, secretList, "my-api-key")
}

func TestDescribeOneSecret(t *testing.T) {
	clientSet := fake.NewSimpleClientset()
	cmd.Namespace = "default"

	for _, s := range testdata.GenerateSecrets {
		_, err := clientSet.CoreV1().Secrets(cmd.Namespace).Create(&s)
		if err != nil {
			assert.Equal(t, nil, err)
		}
	}

	secretDescription, err := cmd.GetSecretDescription(clientSet, "my-credentials")
	if err != nil {
		assert.Equal(t, nil, err)
	}
	assert.Equal(t, "my-credentials", secretDescription.Name)
	assert.Equal(t, []byte("foo"), secretDescription.Data["USERNAME"])
	assert.Equal(t, []byte("bar"), secretDescription.Data["PASSWORD"])
}

func TestDeleteSecret(t *testing.T) {
	clientSet := fake.NewSimpleClientset()
	cmd.Namespace = "default"

	for _, s := range testdata.GenerateSecrets {
		_, err := clientSet.CoreV1().Secrets(cmd.Namespace).Create(&s)
		if err != nil {
			assert.Equal(t, nil, err)
		}
	}

	err := cmd.DeleteSecret(clientSet, "my-credentials")
	if err != nil {
		assert.Equal(t, nil, err)
	}
	secretList, err := clientSet.CoreV1().Secrets(cmd.Namespace).List(metav1.ListOptions{})
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
	cmd.Namespace = "default"

	for _, s := range testdata.GenerateSecrets {
		_, err := clientSet.CoreV1().Secrets(cmd.Namespace).Create(&s)
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

	err := cmd.PatchSecret(clientSet, &secretName, &removeData, &updateData)
	if err != nil {
		assert.Equal(t, nil, err)
	}

	updatedSecret, err := clientSet.CoreV1().Secrets(cmd.Namespace).Get("my-credentials", metav1.GetOptions{})
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
	cmd.Namespace = "default"

	secretName := "new-app-secret"
	container := "default"
	data := []string{
		"USERNAME=foo",
		"PASSWORD=bar",
	}
	err := cmd.CreateSecret(clientSet, &secretName, &container, data)
	if err != nil {
		assert.Equal(t, nil, err)
	}
	createdSecret, err := clientSet.CoreV1().Secrets(cmd.Namespace).Get("new-app-secret-default", metav1.GetOptions{})
	if err != nil {
		assert.Equal(t, nil, err)
	}
	assert.Equal(t, []byte("foo"), createdSecret.Data["USERNAME"])
	assert.Equal(t, []byte("bar"), createdSecret.Data["PASSWORD"])
}
