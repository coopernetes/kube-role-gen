# Generate a Kubernetes RBAC Role

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
Generate a Kubernetes role with every available resource type on a cluster.
Arguments:

  -n,--name - specify the name of the emitted Role. Default is 'foo-clusterrole'
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
  - ''
  resources:
  - bindings
  - pods/binding
  - pods/eviction
  verbs:
  - create
- apiGroups:
  - ''
  resources:
  - componentstatuses
  verbs:
...
```

You can also redirect the output to a file and create your new roles from the generated manifest as a starting point:

```bash
$ kube-role-gen > foo-clusterrole.yaml
2020/02/07 22:42:54 Group: v1
2020/02/07 22:42:54 Resource: bindings - Verbs: [create]
2020/02/07 22:42:54 Resource: componentstatuses - Verbs: [get list]
2020/02/07 22:42:54 Resource: configmaps - Verbs: [create delete deletecollection get list patch update watch]
2020/02/07 22:42:54 Resource: endpoints - Verbs: [create delete deletecollection get list patch update watch]
2020/02/07 22:42:54 Resource: events - Verbs: [create delete deletecollection get list patch update watch]
2020/02/07 22:42:54 Resource: limitranges - Verbs: [create delete deletecollection get list patch update watch]
2020/02/07 22:42:54 Resource: namespaces - Verbs: [create delete get list patch update watch]
2020/02/07 22:42:54 Resource: namespaces/finalize - Verbs: [update]
2020/02/07 22:42:54 Resource: namespaces/status - Verbs: [get patch update]
...

$ kubeval foo-clusterrole.yaml
PASS - foo-clusterrole.yaml contains a valid ClusterRole

$ kubectl apply -f foo-clusterrole.yaml
clusterrole.rbac.authorization.k8s.io/foo-clusterrole created
```
