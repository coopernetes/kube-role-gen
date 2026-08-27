package main

import rego.v1

deny contains msg if {
	input.rules[i].apiGroups[_] == "helloworld.io"
	not valid_crd(input.rules[i].resources)
	msg := "Must contain helloworld.io customresource group"
}

valid_crd(resources) if {
	startswith(resources[_], "helloworlds")
}
