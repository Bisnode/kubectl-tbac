package util

import (
	"fmt"
	"log"
	"os"
	"regexp"

	login "github.com/Bisnode/kubectl-login/util"

	"github.com/pkg/errors"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

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

// CreateClientSet returns a kubernetes clientSet.
func CreateClientSet(ctx *string) (*kubernetes.Clientset, error) {
	configOverrides := &clientcmd.ConfigOverrides{}
	if *ctx != "" {
		configOverrides.CurrentContext = *ctx
	}
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules,
		&clientcmd.ConfigOverrides{CurrentContext: *ctx})
	config, err := kubeConfig.ClientConfig()
	if err != nil {
		return nil, errors.Wrap(err, "failed to load the kube config")
	}

	clientSet, err := kubernetes.NewForConfig(config)
	return clientSet, errors.Wrap(err, "cannot initialize a kubernetes client with loaded configuration")
}

// WhoAmI parses the jwt and looking for groups that it has.
// It matches prefix using matchPrefix and trims away prefix using trimPrefix.
func WhoAmI(matchPrefix, trimPrefix, ctx *string) (teams []string) {
	clientCfg, err := clientcmd.NewDefaultClientConfigLoadingRules().Load()
	if err != nil {
		log.Fatal("Failed to get default config")
	}

	// Set current context if --context is passed.
	if *ctx != "" && *ctx != clientCfg.CurrentContext {
		clientCfg = login.LoadConfigFromContext(*ctx)
		clientCfg.CurrentContext = *ctx
	}

	if clientCfg.CurrentContext == "" {
		fmt.Println("No current-context set - run 'kubectl login --init' to initialize context")
		os.Exit(1)
	}

	rawToken := currentToken(clientCfg)
	if rawToken == "" {
		return []string{}
	}

	claims := login.JwtToIdentityClaims(rawToken)

	return login.ExtractTeams(claims)
}

func currentToken(clientCfg *api.Config) string {
	if clientCfg.CurrentContext == "" {
		fmt.Println("No current-context set - run 'kubectl login --init' to initialize context")
		os.Exit(1)
	}
	// Note that absence of a token is not an error here but an empty string is returned
	return login.ReadToken(clientCfg.CurrentContext)
}
