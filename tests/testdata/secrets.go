package testdata

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GenerateSecrets is a set of secret definitions.
var GenerateSecrets = []v1.Secret{
	v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-credentials",
			Namespace: "default",
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
		},
		Data: map[string][]byte{
			"URL": []byte("github.com"),
			"KEY": []byte("key"),
		},
	},
}
