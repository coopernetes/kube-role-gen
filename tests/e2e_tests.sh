#!/usr/bin/env bash
set -euo pipefail
IFS=$'\n\t'

kube-role-gen | kubeval -
kube-role-gen | kubectl apply --validate -f -
kube-role-gen | conftest test --policy tests/gh-11.rego -

# https://github.com/coopernetes/kube-role-gen/issues/8
kube-role-gen -json | python -m json.tool 2>&1 > /dev/null

# https://github.com/coopernetes/kube-role-gen/issues/14
if [ -f "$HOME/.kube/config" ]; then
    cp $HOME/.kube/config /tmp/test-kubecfg
    KUBECONFIG=/tmp/test-kubecfg kube-role-gen | kubeval -
    kube-role-gen -kubeconfig /tmp/test-kubecfg | kubeval -
fi

kubectl apply -f tests/crd.yaml
kubectl apply -f https://raw.githubusercontent.com/BuddhiWathsala/helloworld-k8s-operator/v0.4.0/deploy/service_account.yaml
kubectl apply -f https://raw.githubusercontent.com/BuddhiWathsala/helloworld-k8s-operator/v0.4.0/deploy/role.yaml
kubectl apply -f https://raw.githubusercontent.com/BuddhiWathsala/helloworld-k8s-operator/v0.4.0/deploy/role_binding.yaml

kubectl wait --for condition=established --timeout=60s crd/helloworlds.helloworld.io

kubectl apply -f https://raw.githubusercontent.com/BuddhiWathsala/helloworld-k8s-operator/v0.4.0/deploy/operator.yaml
kube-role-gen | conftest test --policy tests/gh-7.rego -
