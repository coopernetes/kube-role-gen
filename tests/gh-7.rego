package main

deny[msg] {
        input.rules[i].apiGroups[_] == "helloworld.io"
        not valid_crd(input.rules[i].resources)
        msg := "Must contain helloworld.io customresource group"
}

valid_crd(resources) {
	startswith(resources[_], "helloworlds")
}
