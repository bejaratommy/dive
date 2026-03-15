# Contributing to dive

Thank you for your interest in contributing to dive! This document covers how to get the project building and running locally, and how to submit your changes.

## Prerequisites

- [Go](https://go.dev/dl/) 1.24 or later
- [Docker](https://docs.docker.com/get-docker/) (required for integration/CLI tests)
- `make` or `task` (see [Bootstrapping](#bootstrapping) below)

## Getting the source

```bash
git clone https://github.com/wagoodman/dive.git
cd dive
```

## Bootstrapping

The project uses [task](https://taskfile.dev) as its task runner, managed by [binny](https://github.com/anchore/binny). All tools are pinned and installed locally under `.tool/` — nothing is installed system-wide.

Bootstrap everything in one step:

```bash
make tools
# or equivalently:
make bootstrap
```

From this point you can use either `make <target>` or `.tool/task <target>` interchangeably.

## Building

```bash
make build
```

This produces a local snapshot binary via [goreleaser](https://goreleaser.com). The binary ends up under `./snapshot/`.

For a quick dev build without goreleaser:

```bash
go build -o dive ./cmd/dive
```

## Running the tests

```bash
# Run all validations (static analysis + tests)
make

# Unit tests only
make unit

# CLI/integration tests only (requires Docker)
make cli
```

## Static analysis

```bash
# Run linting and license checks
make static-analysis

# Just lint
make lint

# Check go.mod is tidy
make check-go-mod-tidy
```

Before submitting a PR, make sure `make` (the default target) passes cleanly.

## Project layout

```
cmd/          Entry points (CLI wiring via cobra)
dive/         Core library: image analysis, file-tree diffing
  filetree/   File-tree data structures and operations
  image/      Image fetch/parsing for Docker and Podman
internal/     Internal utilities
runtime/      TUI, CI mode, export, and test-CLI
  ci/         CI pass/fail evaluation logic
  export/     JSON export
  ui/         Terminal UI (gocui-based)
```

## Submitting changes

1. Fork the repository and create a branch off `main`.
2. Make your changes. Keep commits focused and the diff minimal.
3. Run `make` and confirm everything passes.
4. Open a pull request against `wagoodman/dive:main`. Fill in a brief description of what you changed and why, and link any relevant issues.

## Code style

- Follow standard Go conventions (`gofmt`, `goimports`).
- The project uses [golangci-lint](https://golangci-lint.run/) with the configuration in `.golangci.yaml`. Running `make lint` locally will catch most issues before CI does.
- Match the style of the surrounding code when in doubt.

## License

By contributing you agree that your changes will be licensed under the [MIT License](LICENSE).
