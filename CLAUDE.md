# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project

LoseNotion is a terminal todo-list app built so the author doesn't have to open Notion for small notes. The author is also using it to learn Go, so prefer idiomatic, readable Go and explain non-obvious language or Bubble Tea concepts when making changes.

It is at an early stage: todos can be listed and added, and are held in memory only (no persistence yet).

## Commands

```bash
go run .                 # run the TUI (Ctrl+C to quit)
go build -o losenotion . # build a binary
go vet ./...             # static checks
gofmt -l .               # list unformatted files (use gofmt -w . to fix)
go test ./...            # run tests
go test -run 'TestAdd/ADD.3' ./...  # run one subtest (spaces in subtest names become _)
go mod tidy              # run after adding/removing imports
```

## Architecture

Single `main` package at the repo root, built on [Bubble Tea](https://github.com/charmbracelet/bubbletea) (v1, imported as `tea`) using the Elm architecture:

- `main.go` constructs a `tea.Program` with the root model (`todoList`) and runs it. `simple_page.go` is an earlier experiment and is no longer used.
- Models implement `tea.Model` (`Init`, `Update`, `View`) with value receivers and a `newXxx(...)` constructor. `Update` returns a new copy of the model rather than mutating it. Copies share slice backing arrays, so tests must always continue from the model `Update` returns.
- `todoList` (`todo_list.go`) switches views with a `mode` field instead of separate models, so all state stays in one model. `Update` handles `ctrl+c` first, then hands off to a per-mode `updateXxx` method.
- Quitting is handled by the model itself. Any new root model must keep handling `ctrl+c`, or the program can't be exited.
- Tests drive the model by sending `tea.KeyMsg`s through the `press()` helper in `todo_list_test.go`, with no terminal.

## Specs (acai.sh)

Behavior is specified in `features/lose-notion/*.feature.yaml` using acai.sh conventions; run `npx @acai.sh/cli skill` for the full process. Update the spec before changing behavior, never renumber requirements, and reference full ACIDs (e.g. `todo-list.ADD.3`) in `//` comments and in subtest names (`t.Run("todo-list.ADD.3 ...")`, since Go test function names can't contain dots).

The acai server is self-hosted at http://localhost:4000. The CLI reads config only from environment variables and silently falls back to the public app.acai.sh if `ACAI_API_BASE_URL` is unset, so load `.env` in the same command as every call:

```bash
set -a && . ./.env && set +a && npx @acai.sh/cli push --all
```
