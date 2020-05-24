module github.com/Bisnode/kubectl-tbac

go 1.14

require (
	cloud.google.com/go v0.38.0 // indirect
	github.com/Azure/go-autorest v11.1.2+incompatible // indirect
	github.com/Azure/go-autorest/autorest v0.9.0 // indirect
	github.com/Bisnode/kubectl-login v1.1.0
	github.com/golang/protobuf v1.3.5 // indirect
	github.com/gophercloud/gophercloud v0.1.0 // indirect
	github.com/gregjones/httpcache v0.0.0-20180305231024-9cad4c3443a7 // indirect
	github.com/imdario/mergo v0.3.9 // indirect
	github.com/json-iterator/go v1.1.9 // indirect
	github.com/peterbourgon/diskv v2.0.1+incompatible // indirect
	github.com/pkg/errors v0.9.1
	github.com/skratchdot/open-golang v0.0.0-20190402232053-79abb63cd66e // indirect
	github.com/spf13/cobra v0.0.7
	github.com/stretchr/testify v1.4.0
	golang.org/x/crypto v0.0.0-20200323165209-0ec3e9974c59 // indirect
	golang.org/x/net v0.0.0-20200324143707-d3edc9973b7e // indirect
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d // indirect
	golang.org/x/sys v0.0.0-20200327173247-9dae0f8f5775 // indirect
	golang.org/x/time v0.0.0-20191024005414-555d28b269f0 // indirect
	google.golang.org/appengine v1.6.5 // indirect
	k8s.io/api v0.15.11
	k8s.io/apimachinery v0.17.0
	k8s.io/client-go v11.0.0+incompatible
	k8s.io/utils v0.0.0-20200327001022-6496210b90e8 // indirect
	sigs.k8s.io/structured-merge-diff v0.0.0-20190525122527-15d366b2352e // indirect
)

replace k8s.io/client-go => k8s.io/client-go v0.15.11
