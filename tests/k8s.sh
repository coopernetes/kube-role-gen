#!/usr/bin/env bash
set -euo pipefail
IFS=$'\n\t'

kube-role-gen | kubeval -
kube-role-gen | kubectl apply --validate -f -
kube-role-gen | conftest test --policy tests/gh-11.rego -

# https://github.com/coopernetes/kube-role-gen/issues/14
if [ -f "$HOME/.kube/config" ]; then
    mv $HOME/.kube/config $HOME/.kube/config.bak
    KUBECONFIG=/tmp/test-kubecfg kube-role-gen
    kube-role-gen -kubeconfig $HOME/.kube/config.bak
fi
