# Kubernetes TBAC manager
## (!) This is a work in progress.


This is a kubectl plugin that is built to simplify creating resources in a team's namespace.
With team based access control is should be possible to manage secrets and configmaps in the development teams own namespace.

Use of this plugin requires no knowledge of Kubernetes manifests and namespaces. A team member can be sure that the resource is created with the correct labels, using the correct requirements for the Kubernetes API and that it ends up in the right place in Kubernetes.

# Install
## Ubuntu
Download or build the binary locally.

Make it executable: `chmod +x kubectl-tbac`

Put it in your `$PATH`.

# Build locally

`go build .`