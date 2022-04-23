# kube-role-gen - Create a complete Kubernetes RBAC Role

[![Go](https://github.com/coopernetes/kube-role-gen/workflows/Go/badge.svg)](https://github.com/coopernetes/kube-role-gen/actions?query=workflow%3AGo)
[![Go Report Card](https://goreportcard.com/badge/github.com/coopernetes/kube-role-gen)](https://goreportcard.com/report/github.com/coopernetes/kube-role-gen)

`kube-role-gen` is a command-line utility that will generate a Kubernetes ClusterRole that contains every resource available on a connected cluster, including sub-resources & custom resources. All rules will be grouped by their unique API group & supported verbs combinations so a granular ClusterRole or Role may be created.

This is useful when you want to define a role with broad set of permissions while explicitly excluding a small subset of them. An example might be a cluster administrator who should have no access to namespace Secrets, ServiceAccounts or RBAC Roles/Bindings. Until Kubernetes supports [substraction via role aggregation](https://github.com/kubernetes/kubernetes/issues/70387), this tool is useful as a starting point to generate roles programmatically.

Another neat feature of `kube-role-gen` is defining a role that relies heavily on sub-resources. Sub-resources such as `pod/exec` do not show up in any static list such as `kubectl api-resources -o wide` and must be discovered by interacting with the Kubernetes API directly. See this [stackoverflow answer for additional details](https://stackoverflow.com/a/51289417).

This utility was inspired by [this original bash implementation](https://stackoverflow.com/a/57892189).

## Install

Download the latest [release](https://github.com/coopernetes/kube-role-gen/releases):

```bash
curl -LO https://github.com/coopernetes/kube-role-gen/releases/download/v0.0.4/kube-role-gen_0.0.4_Linux_x86_64.tar.gz
tar xf kube-role-gen_0.0.4_Linux_x86_64.tar.gz
mv kube-role-gen /usr/local/bin/
```

You can also install as a Go module. Ensure you have `$GOPATH/bin` on your `$PATH`:

```bash
PATH="$(go env GOPATH)/bin:$PATH"
GO111MODULE="on" go get github.com/coopernetes/kube-role-gen
```

## Usage

```bash
$ kube-role-gen -h
Usage of kube-role-gen:
  -json
        Generate JSON output. If unset, will default to YAML.
  -kubeconfig string
        absolute path to the kubeconfig file. If set, this will override the default behavior and ignore KUBECONFIG environment variable and/or $HOME/.kube/config file location.
  -name string
        Override the name of the ClusterRole resource that is generated (default "foo-clusterrole")
  -pretty
        Enable human-readable JSON output. This flag is ignored for YAML (always pretty-prints).
  -v    Enable verbose logging.
```

The resulting `ClusterRole` resource will be printed to stdout in YAML format.

```bash
$ kube-role-gen
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  creationTimestamp: null
  name: foo-clusterrole
rules:
- apiGroups:
  - ""
  resources:
  - bindings
  - pods/binding
  - pods/eviction
  - serviceaccounts/token
  verbs:
  - create
- apiGroups:
  - ""
  resources:
  - componentstatuses
  verbs:
  - get
  - list
- apiGroups:
  - ""
  resources:
  - configmaps
  - endpoints
  - events
  - limitranges
  - nodes
  - persistentvolumeclaims
  - persistentvolumes
  - pods
  - podtemplates
  - replicationcontrollers
  - resourcequotas
  - secrets
  - serviceaccounts
  verbs:
  - create
  - delete
  - deletecollection
  - get
  - list
  - patch
  - update
  - watch
...
```

You can also redirect the output to a file and create your new roles from the generated manifest as a starting point:

```bash
$ kube-role-gen > foo-clusterrole.yaml

$ kubeval foo-clusterrole.yaml
PASS - foo-clusterrole.yaml contains a valid ClusterRole

$ kubectl apply -f foo-clusterrole.yaml
clusterrole.rbac.authorization.k8s.io/foo-clusterrole created
```
