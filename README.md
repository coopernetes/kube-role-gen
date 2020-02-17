# kube-role-gen - Create a complete Kubernetes RBAC Role

![Go](https://github.com/coopernetes/kube-role-gen/workflows/Go/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/coopernetes/kube-role-gen)](https://goreportcard.com/report/github.com/coopernetes/kube-role-gen)

This binary will generate a valid Kubernetes RBAC role that contains every resource available on a cluster, including subresources. It will walk the API of the connected cluster and gather all available resources. All discovered resources will be grouped by their unique API group & supported verbs combinations so a complete & granular RBAC role may be created.

This is useful when you want to define a broad RBAC role that has access to _most_ objects but is disallowed from viewing a subset. Until Kubernetes supports [substraction via role aggregation](https://github.com/kubernetes/kubernetes/issues/70387), this script is useful as a starting point.

Another use case is defining a role that relies heavily on sub-resources. Sub-resources such as `pod/exec` do not show up in any static list such as `kubectl api-resources -o wide` and must be discovered by walking the Kubernetes API. See this [stackoverflow answer for additional details](https://stackoverflow.com/a/51289417).

This utility was inspired by [this original bash implementation](https://stackoverflow.com/a/57892189).

## Install

Ensure that your GOPATH is included in your path:

```bash
PATH="$PATH:$(go env GOPATH)/bin"
```

```bash
GO111MODULE="on" go get github.com/coopernetes/kube-role-gen
```

## Usage

```bash
$ kube-role-gen -h
Usage of kube-role-gen:
  -kubeconfig string
        (optional) absolute path to the kubeconfig file (default "/home/tom/.kube/config")
  -name string
        Override the name of the ClusterRole resource that is generated (default "foo-clusterrole")
  -v    Enable verbose logging
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
  - componentstatuses
  - configmaps
  - endpoints
  - events
  - limitranges
  - namespaces
  - namespaces/finalize
  - namespaces/status
  - nodes
  - nodes/proxy
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
