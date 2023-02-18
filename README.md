# kube-role-gen - Create a complete Kubernetes RBAC Role

[![Go](https://github.com/coopernetes/kube-role-gen/workflows/Go/badge.svg)](https://github.com/coopernetes/kube-role-gen/actions?query=workflow%3AGo)
[![Go Report Card](https://goreportcard.com/badge/github.com/coopernetes/kube-role-gen)](https://goreportcard.com/report/github.com/coopernetes/kube-role-gen)

`kube-role-gen` is a command-line utility that will generate a Kubernetes
ClusterRole that contains every resource available on a connected cluster,
including sub-resources & custom resources. Each entry in the ClusterRole
rules are grouped by API group and the combination of unique resource type &
supported verbs. This is different from something like
`kubectl create role ... -o yaml --dry-run=client`, which groups resources
together even if they all don't support the same verb (ie. `pods/exec` listed
with the `patch` verb).

## Why create this?
* In secure environments, even cluster admins shouldn't have access to
  everything. Access to resources such as namespace creation/delete,
  rolebindings, etc. should be reserved for cluster management tools, pipelines
  or scripts.
* Kubernetes will likely never support [role aggregation via subtraction](https://github.com/kubernetes/kubernetes/issues/70387).
* Sub-resources such as `pods/exec` are not accessible via any normal `kubectl`
  output (with the only exception being `kubectl --raw`). It must be queried
  using Kubernetes API discovery via a client.
* I didn't want to maintain the [original bash script](https://stackoverflow.com/a/57892189)
  to do the same thing. Props to [Vit on stackoverflow](https://stackoverflow.com/users/9929015/vit)
  for providing the idea for this utility.
* It's my own excuse to learn Go for something I need at work.

_Alternatives_:
- Use privileged access management for any elevated permissions inside Kubernetes.
- Use a tool such as [audit2rbac](https://github.com/liggitt/audit2rbac) to
  generate a least-privilege role based on what your cluster users are actually deploying

## Install

Download the latest [release](https://github.com/coopernetes/kube-role-gen/releases):

```bash
curl -LO https://github.com/coopernetes/kube-role-gen/releases/download/v0.0.5/kube-role-gen_0.0.5_Linux_x86_64.tar.gz
tar xf kube-role-gen_0.0.5_Linux_x86_64.tar.gz
mv kube-role-gen /usr/local/bin/
```

You can also install as a Go module.

```bash
go install github.com/coopernetes/kube-role-gen
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


## Manipulation / Post-Processing
This utility doesn't provide any post-processing out of the box. However, you can use tools such as `jq` to
chain the output of kube-role-gen and manipulate it as you see fit. 

Here's a few common "recipes" for manipulating the role that is generated:

_No delete access_
```shell
kube-role-gen -json | jq 'del(.rules[].verbs[] |           
select((. == "delete") or (. == "deletecollection")))'
```

_Read-only access to all resources_
```shell
kube-role-gen -json | jq 'del(.rules[].verbs[] |           
select((. == "create") or (. == "delete") or (. == "deletecollection") or (. == "patch") or (. == "update")))'
```

_Exclude a specific API group_
```shell
kube-role-gen -json | jq 'del(.rules[] | select(.apiGroups[] | contains("flowcontrol.apiserver.k8s.io")))' 
```

_Exclude multiple API groups_
```shell
kube-role-gen -json | jq 'del(.rules[] | select(.apiGroups[] | contains("scheduling.k8s.io") or contains("flowcontrol.apiserver.k8s.io") or contains("node.k8s.io")))'
```
