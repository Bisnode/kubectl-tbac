module github.com/Bisnode/kubectl-tbac

go 1.14

require (
	github.com/Bisnode/kubectl-login v1.2.0
	github.com/pkg/errors v0.9.1
	github.com/spf13/cobra v1.0.0
	github.com/stretchr/testify v1.6.1
	k8s.io/api v0.15.11
	k8s.io/apimachinery v0.19.1
	k8s.io/client-go v11.0.0+incompatible
)

replace k8s.io/client-go => k8s.io/client-go v0.15.11

replace k8s.io/apimachinery => k8s.io/apimachinery v0.17.0
