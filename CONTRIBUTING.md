# Contributing

Thanks for contributing to `kube-role-gen`! This is a small Go CLI; local
development is intentionally simple and mirrors what CI does.

## Prerequisites

- [mise](https://mise.jdx.dev/) to manage tool versions. Everything the project
  needs (Go, kind, kubectl, conftest, kubeconform) is declared in `mise.toml`:

  ```bash
  mise install
  ```

- A container runtime for the end-to-end tests, since they spin up a
  [kind](https://kind.sigs.k8s.io/) cluster. Docker works out of the box; with
  Podman, export `KIND_EXPERIMENTAL_PROVIDER=podman`.

> CI does **not** use mise — it installs the same tools via first-party GitHub
> Actions. mise is just a convenience for local development.

## Build, test, lint

Use the `Makefile` targets:

```bash
make test      # go test ./...
make build     # build the ./kube-role-gen binary
make install   # go install into $(go env GOPATH)/bin
make check     # gofmt check + go vet
make fmt       # gofmt -w
```

## End-to-end tests

The e2e suite (`tests/e2e_tests.sh`, run via `make e2e`) generates a ClusterRole
against a live cluster and validates it with `kubeconform`, `kubectl apply`, and
`conftest` policies. You need a running cluster reachable via your current
kubeconfig context:

```bash
# create a throwaway cluster (add KIND_EXPERIMENTAL_PROVIDER=podman if needed)
kind create cluster --name krg-e2e

# make sure the binary is on your PATH, then run the suite
make install
export PATH="$(go env GOPATH)/bin:$PATH"
make e2e

# tear it down when finished
kind delete cluster --name krg-e2e
```

## What CI runs

`.github/workflows/build.yml` runs on every push and pull request and is the
source of truth. It mirrors the above: `make test`, `make install`, then creates
a kind cluster and runs `make e2e`. If those pass locally, they should pass in
CI.

## Pull requests

- Keep changes focused and run `make check` and `make test` before opening a PR.
- Add or update tests for behaviour changes (`pkg/k8s/rbac_test.go` for role
  generation logic, the `tests/` policies for e2e expectations).
