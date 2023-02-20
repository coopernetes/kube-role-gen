#!/usr/bin/env bash
set -euo pipefail
IFS=$'\n\t'

kube-role-gen | kubeconform -summary
kube-role-gen | kubectl apply --validate -f -
kube-role-gen | conftest test --policy tests/gh-11.rego -

# https://github.com/coopernetes/kube-role-gen/issues/8
if command -v python &> /dev/null
then
  kube-role-gen -json | python -m json.tool > /dev/null 2>&1
else
  kube-role-gen -json | python3 -m json.tool > /dev/null 2>&1
fi

# https://github.com/coopernetes/kube-role-gen/issues/14
if [ -f "$HOME/.kube/config" ]; then
    cp $HOME/.kube/config /tmp/test-kubecfg
    KUBECONFIG=/tmp/test-kubecfg kube-role-gen | kubeconform -summary
    kube-role-gen -kubeconfig /tmp/test-kubecfg | kubeconform -summary
fi

kubectl apply --validate=false -f tests/crd.yaml
kube-role-gen | conftest test --policy tests/gh-7.rego -
