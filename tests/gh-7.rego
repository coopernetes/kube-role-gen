package main

deny[msg] {
	input.rules[i].apiGroups[_] == "helloworld.io"
	msg := "Must contain helloworld customresource"
}

