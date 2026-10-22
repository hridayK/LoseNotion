# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project

LoseNotion is a terminal todo-list app built so the author doesn't have to open Notion for small notes. The author is also using it to learn Go, so prefer idiomatic, readable Go and explain non-obvious language or Bubble Tea concepts when making changes.

## Commands

```bash
go run .                 # run the TUI
go run . reminders -interval 2s  # run the reminder notifier in the foreground
go build -o losenotion . # build a binary
go vet ./...             # static checks
gofmt -l .               # list unformatted files (use gofmt -w . to fix)
go test ./...            # run all tests
go test -run 'TestAdd/ADD.3' ./internal/tui/  # run one subtest (spaces in subtest names become _)
go mod tidy              # run after adding/removing imports
```

Set `LOSENOTION_DATA_DIR` to a scratch directory when running the app manually, so the user's real data in `~/.local/share/losenotion` isn't touched.

## Architecture

`main.go` only wires things up: no arguments runs the TUI, `reminders` runs the notifier. Everything else is in `internal/`, layered so lower packages know nothing about higher ones:

- `todo`: the `List` type (active todos plus deleted ones) and every operation on it. No I/O, no UI. Methods ignore out-of-range indexes so callers don't need to guard.
- `store`: loads and saves the list as `todos.json` (atomic temp file + rename), plus the notifier's `reminders-last-check` timestamp.
- `reminder`: `Checker` finds reminders in the window (last check, now] and sends one notification per todo via `SystemNotifier` (`notify-send`/`osascript`). It never writes `todos.json`, only its own last-check file, so it can run alongside the TUI without write conflicts.
- `tui`: one Bubble Tea `Model`. Screens are a `mode` field rather than separate models, with one file per screen (`list.go`, `trash.go`, `reminders.go`), and `prompt.go` handles the three text-input modes. `Update` handles `ctrl+c` first, then dispatches by mode. Every change calls `persist()`, which saves straight away through the injected `save` func.

`tui.New` takes `save` and `now` as parameters, and `reminder.Checker` takes a `Store` interface, a `Notifier` and a `Now` func, so tests run without a terminal, real files, or the real clock. Models use value receivers; copies share slice backing arrays, so tests always continue from the model `Update` returns (the `press()` helper in `internal/tui/helpers_test.go` does this).

## Specs (acai.sh)

Behavior is specified in `features/lose-notion/*.feature.yaml` (`todo-list`, `persistence`, `reminders`) using acai.sh conventions; run `npx @acai.sh/cli skill` for the full process. Update the spec before changing behavior, never renumber requirements (mark them `deprecated: true` instead), and reference full ACIDs (e.g. `todo-list.ADD.3`) in `//` comments and in subtest names (`t.Run("todo-list.ADD.3 ...")`, since Go test function names can't contain dots).

The acai server is self-hosted at http://localhost:4000. The CLI reads config only from environment variables and silently falls back to the public app.acai.sh if `ACAI_API_BASE_URL` is unset, so load `.env` in the same command as every call:

```bash
set -a && . ./.env && set +a && npx @acai.sh/cli push --all
```
