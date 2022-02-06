#!/usr/bin/env bash
KUBE_VERSION=$KUBE_VERSION
if [ -z $KUBE_VERSION ]; then
    KUBE_VERSION="$(curl -L -s https://dl.k8s.io/release/stable.txt)"
fi

curl -LO "https://dl.k8s.io/release/${KUBE_VERSION}/bin/linux/amd64/kubectl"

curl -LO "https://dl.k8s.io/${KUBE_VERSION}/bin/linux/amd64/kubectl.sha256"

echo "$(<kubectl.sha256)  kubectl" | sha256sum --check

sudo install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl

rm -f kubectl
rm -f kubectl.sha256

minikube delete
minikube start --kubernetes-version=$KUBE_VERSION
