# AGENTS.md

## Purpose

Flashcards is a Go CLI that generates and reviews local flashcards with Ollama, Cobra, Bubble Tea, and SQLite.

## Repo shape

- `cmd/flashcards`: CLI entrypoint
- `internal/commands`: Cobra commands and app wiring
- `internal/tui`: Bubble Tea models, layout, theme, components
- `internal/store`: SQLite persistence
- `internal/config`: env-driven config
- `internal/ollama`: Ollama integration and response parsing
- `internal/security`: path and input validation

## Make changes safely

- prefer small diffs
- preserve package boundaries
- keep imports on module path `github.com/jae-labs/flashcards/...`
- add tests for logic changes
- do not mix unrelated cleanup into feature work

## Validate before done

```bash
make fmt
make test
make lint
```

## Git hooks

This repo uses Lefthook via `lefthook.yml`.

```bash
lefthook install
```

Current hooks:
- `pre-commit`: `gofmt`, `golangci-lint`
- `pre-push`: `go test -v -race ./...`

## Notes for agents

- read `Makefile` before adding new dev commands
- prefer updating existing patterns over adding parallel abstractions
- flag lint/test debt you did not introduce
- if a command fails, report exact failure and next fix path
