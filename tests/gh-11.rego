package main

import rego.v1

deny contains msg if {
	input.rules[i].apiGroups[_] == "batch"
	not valid_batch(input.rules[i].resources)
	msg := "Must contain all batch resources"
}

valid_batch(resources) if {
	startswith(resources[_], "cronjobs")
	startswith(resources[_], "jobs")
}
