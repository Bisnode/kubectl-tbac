module github.com/Bisnode/kubectl-tbac

go 1.14

require (
	github.com/Bisnode/kubectl-login v1.1.1
	github.com/golang/protobuf v1.4.2 // indirect
	github.com/imdario/mergo v0.3.9 // indirect
	github.com/pkg/errors v0.9.1
	github.com/spf13/cobra v1.0.0
	github.com/stretchr/testify v1.6.1
	golang.org/x/crypto v0.0.0-20200604202706-70a84ac30bf9 // indirect
	golang.org/x/net v0.0.0-20200602114024-627f9648deb9 // indirect
	golang.org/x/sys v0.0.0-20200615200032-f1bc736245b1 // indirect
	k8s.io/api v0.18.3
	k8s.io/apimachinery v0.18.3
	k8s.io/client-go v11.0.0+incompatible
	k8s.io/utils v0.0.0-20200603063816-c1c6865ac451 // indirect
)

replace k8s.io/client-go => k8s.io/client-go v0.15.11
