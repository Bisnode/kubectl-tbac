package util

import (
	"log"
	"regexp"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

// IdentityClaims ...
type IdentityClaims struct {
	Username string    `json:"email"`
	Groups   *[]string `json:"groups"`
	jwt.StandardClaims
}

// CreateClientSet returns a kubernetes clientset.
func CreateClientSet() (*kubernetes.Clientset, error) {
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules,
		&clientcmd.ConfigOverrides{})
	config, err := kubeConfig.ClientConfig()
	if err != nil {
		return nil, errors.Wrap(err, "failed to load the kube config")
	}

	clientset, err := kubernetes.NewForConfig(config)
	return clientset, errors.Wrap(err, "cannot initialize a kubernetes client with loaded configuration")
}

// WhoAmI ...
func WhoAmI() (teams []string) {
	parser := &jwt.Parser{}
	claims := &IdentityClaims{}

	clientCfg, err := clientcmd.NewDefaultClientConfigLoadingRules().Load()
	if err != nil {
		log.Fatal("Failed to get default config")
	}
	rawToken := currentToken(clientCfg)
	if rawToken == "" {
		return []string{}
	}

	_, _, err = parser.ParseUnverified(rawToken, claims)
	if err != nil {
		log.Fatalf("Failed parsing token: %v", rawToken)
	}

	for _, group := range *claims.Groups {
		group := strings.ToLower(group)
		if !strings.HasPrefix(group, "sec-tbac-team-") {
			continue
		}
		teams = append(teams, strings.TrimPrefix(group, "sec-tbac-"))
	}

	return teams
}

func currentToken(clientCfg *api.Config) string {
	if clientCfg.CurrentContext == "" {
		return ""
	}

	if clientCfg.AuthInfos[clientCfg.CurrentContext].Token == "" {
		return ""
	}

	return clientCfg.AuthInfos[clientCfg.CurrentContext].Token
}

// AssembleInputData is meant to parse data key value pairs
// coming from the command line and to be put in the
// data field in a secret or configmap.
func AssembleInputData(data []string) map[string][]byte {
	dataMap := make(map[string][]byte)
	for _, kvp := range data {
		kv := regexp.MustCompile(`=`).Split(kvp, 2)
		dataMap[kv[0]] = []byte(kv[1])
	}
	return dataMap
}
