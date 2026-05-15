# AGENTS.md

This file provides guidance to coding agents (e.g. Claude Code, claude.ai/code) when working with code in this repository.

## Repository purpose

Go module `kubeops.dev/falco-ui-server` — an aggregated Kubernetes API server that exposes [Falco](https://falco.org/) runtime-security events as a native Kubernetes resource. [Falcosidekick](https://github.com/falcosecurity/falcosidekick) forwards Falco events to this server's webhook, which stores and re-presents them under `falco.appscode.com/v1alpha1` so any kubectl-compatible UI can list/query them.

The produced binary is `falco-ui-server`. Long-running aggregated apiserver.

## Architecture

- `cmd/falco-ui-server/main.go` — entry point.
- `pkg/cmds/` — Cobra command tree.
- `pkg/apiserver/` — aggregated API server config and lifecycle.
- `pkg/registry/falco/` — `rest.Storage` implementation for `FalcoEvent` (the Kubernetes object exposed via the apiserver).
- `apis/falco/v1alpha1/` — Kubebuilder API types:
  - `falco_event_types.go` — `FalcoEvent`.
  - `register.go`, `doc.go`, generated `zz_generated.{deepcopy,defaults,conversion}.go` and `openapi_generated.go`.
- `client/` — generated clientset.
- `crds/` — generated CRD YAMLs.
- `pkg/falcosidekick/`:
  - `handler.go` — HTTP handler that accepts falcosidekick's `POST /events` payload.
  - `reconciler.go` — turns ingested events into `FalcoEvent` objects.
  - `config.go`, `types/` — shared config and payload types.
  - `metricshandler/` — Prometheus metrics endpoint.
- `pkg/cleaner/` — retention/cleanup of old events.
- `Dockerfile.in` (PROD, distroless), `Dockerfile.dbg` (debian), `Dockerfile.ubi` (Red Hat certified) — three image variants.
- `PROJECT` — Kubebuilder metadata.
- `hack/`, `Makefile` — AppsCode build harness.
- `vendor/` — checked-in deps.

CRD API group is `falco.appscode.com`.

## Common commands

All Make targets run inside `ghcr.io/appscode/golang-dev` — Docker must be running.

- `make ci` — CI pipeline.
- `make build` / `make all-build` — build host or all-platform binaries.
- `make gen` — regenerate clientset + manifests + openapi + conversions/defaults. Run after any change to `apis/falco/v1alpha1/*_types.go`.
- `make manifests` — regenerate CRDs only.
- `make clientset` — regenerate `client/` only.
- `make openapi` — regenerate OpenAPI definitions.
- `make fmt`, `make lint`, `make unit-tests` / `make test` — standard.
- `make verify` — `verify-gen verify-modules`; `go mod tidy && go mod vendor` must leave the tree clean.
- `make container` — build PROD, DBG, and UBI images.
- `make push` — push all three; `make docker-manifest` writes multi-arch manifests; `make release` is the full publish flow.
- `make push-to-kind` / `make deploy-to-kind` — load into Kind and Helm-install.
- `make install` / `make uninstall` / `make purge` — Helm install lifecycle.
- `make add-license` / `make check-license` — manage license headers.

Run a single Go test (requires a local Go toolchain):

```
go test ./pkg/falcosidekick/... -run TestName -v
```

## Conventions

- Module path is `kubeops.dev/falco-ui-server` (vanity URL). Imports must use that.
- License: `LICENSE` (Apache-2.0); new files need the standard AppsCode header (`make add-license`).
- Sign off commits (`git commit -s`); contributions follow the DCO.
- Vendor directory is checked in — `go mod tidy && go mod vendor` must leave the tree clean (enforced by `verify-modules`).
- Do not hand-edit `zz_generated.*.go`, `openapi_generated.go`, anything under `client/`, or `crds/` — change `apis/falco/v1alpha1/*_types.go` and re-run `make gen`.
- This is an **aggregated apiserver**, not a controller-runtime app. Persistence goes through `pkg/registry/falco/`; do not introduce a parallel storage path.
- The falcosidekick → server contract is the HTTP payload at `pkg/falcosidekick/types/`. Don't change that shape without coordinating with deployed falcosidekick releases.
- Three Dockerfiles, one binary — keep `Dockerfile.in`, `Dockerfile.dbg`, and `Dockerfile.ubi` in sync.
