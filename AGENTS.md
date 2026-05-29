# AGENTS.md

This file provides guidance to Codex (Codex.ai/code) when working with code in this repository.

## Context

This directory is currently empty and located within a Go workspace at `D:\GoWork\src\2309A\Codex`. Neighboring sibling directories (1–6) contain Go projects and exercises.

## Go Workspace

The GOPATH is `D:\GoWork`. Source code lives under `D:\GoWork\src\`. When initializing a new module here:

```bash
go mod init <module-path>
```

## Common Go Commands

```bash
# Build
go build ./...

# Run tests
go test ./...

# Run a single test
go test ./path/to/pkg -run TestName

# Lint (if golangci-lint is installed)
golangci-lint run

# Format
go fmt ./...

# Vet
go vet ./...
```

## Notes

- Update this file once the project is initialized with its actual structure, architecture, and commands.
