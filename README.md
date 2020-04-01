# kubectl-tbac

This is a kubectl plugin that is built to simplify creating resources in a team's namespace.

Use of this plugin requires no previous knowledge of Kubernetes manifests and namespaces. A team member can be sure that the resource is created with the correct labels, using the correct requirements for the Kubernetes API and that it ends up in the right place in Kubernetes.
The plugin makes use of the ID token in the kubeconfig to figure out which team the user belong to.

Currently supported resources:
* Secrets

# Install
Download the kubectl-tbac binary and place it in your $PATH.
Once in your $PATH you can start using `kubectl tbac`.

# Build
```
git clone git@github.com:Bisnode/kubectl-tbac.git
cd kubectl-tbac
$GO111MODULE=auto go build
```

# Usage
Some examples of how to manage kubernetes secrets using kubectl-tbac.

Create secret
```
kubectl tbac create secret my-secret --data "USERNAME=foo" --data "PASSWORD=bar"
```

Update secret
```
kubectl tbac create secret my-secret --data "URL=github.com" --data "USERNAME=bar" --remove-data "PASSWORD=bar"
```

List secrets
```
kubectl tbac get secrets
```

Describe one secret
```
kubectl tbac get secret my-secret
```

Delete secret
```
kubectl tbac delete secret my-secret
```

*All commands accepts an --[h]elp flag for more information and examples.*