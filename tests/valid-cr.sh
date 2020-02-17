#!/usr/bin/env bash

go build
./kube-role-gen | kubeval -
