package main

deny[msg] {
	input.rules[i].apiGroups[_] == "batch"
	not valid_batch(input.rules[i].resources)
	msg := "Must contain all batch resources"
}

valid_batch(resources) {
  startswith(resources[_], "cronjobs")
  startswith(resources[_], "jobs")
}
