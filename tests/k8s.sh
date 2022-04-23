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
