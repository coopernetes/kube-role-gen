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
$ python3 generate-role.py
2020-01-03 11:41:54,534 - INFO - Gathering core API resource details
2020-01-03 11:41:54,534 - INFO - Gathering API groups & resource details
2020-01-03 11:41:59,661 - INFO - Resource discovery complete. Found 76 resources in 19 API groups
2020-01-03 11:41:59,661 - INFO - Converting resources to rbac.authorization.k8s.io/v1/Role
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: foo-role
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

You can also redirect the output to a file and create your new Roles from the generated manifest as a starting point:

```bash
$ python3 generate-role.py > foo-role.yaml
2020-01-03 11:42:07,417 - INFO - Gathering core API resource details
2020-01-03 11:42:07,417 - INFO - Gathering API groups & resource details
2020-01-03 11:42:12,676 - INFO - Resource discovery complete. Found 76 resources in 19 API groups
2020-01-03 11:42:12,677 - INFO - Converting resources to rbac.authorization.k8s.io/v1/Role

$ kubeval foo-role.yaml
PASS - foo-role.yaml contains a valid Role

$ kubectl apply -f foo-role.yaml
role.rbac.authorization.k8s.io/foo-role created
```
