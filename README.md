# Generate a Kubernetes RBAC Role

This script will generate a valid Kubernetes RBAC role that contains every resource available on a cluster, including subresources. It will walk the API of the connected cluster and gather all available resources. All discovered resources will be grouped by their unique API group & supported verbs combinations so a complete & granular RBAC role may be created.

This is useful when you want to define a broad RBAC role that has access to _most_ objects but is disallowed from viewing a subset. Until Kubernetes supports [substraction via role aggregation](https://github.com/kubernetes/kubernetes/issues/70387), this script is useful as a starting point.

Another use case is defining a role that relies heavily on sub-resources. Sub-resources such as `pod/exec` do not show up in any static list such as `kubectl api-resources -o wide` and must be discovered by walking the Kubernetes API. See this [stackoverflow answer for additional details](https://stackoverflow.com/a/51289417).

This script was based on [this original bash implementation](https://stackoverflow.com/a/57892189).

## Usage

Requires:

* Python 3.7+

* kubectl

The resulting `Role` resource will be printed to stdout in YAML format.

```bash
./generate-role.py
```

## Validation

```bash
$ ./generate-role.py > foo-role.yaml
$ kubeval foo-role.yaml
PASS - foo-role.yaml is a valid Role
$ kubectl apply -f foo-role.yaml
role/foo-role created
```
