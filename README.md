# LoseNotion
A CLI todo list app because I don't like to open Notion for the smallest of notes. (Also to learn Go)

## Install

```bash
go build -o ~/.local/bin/losenotion .
```

## Use

Run `losenotion` to open your list.

| Key | List | Deleted todos (`t`) | Reminders (`r`) |
|---|---|---|---|
| `↑`/`↓` or `k`/`j` | move | move | move |
| `a` | add a todo | | add a reminder |
| `space` | mark done / not done | | |
| `d` | delete (moves it to deleted todos) | | remove reminder |
| `u` | set due date (empty clears it) | | |
| `r` | open reminders | restore | |
| `x` | | delete forever | |
| `t` | open deleted todos | back | |
| `esc` | | back | back |
| `q` | quit | | |
| `ctrl+c` | quit from anywhere | | |

Dates and times are typed as `2026-10-23 18:00`, or relative to now with `+30m`, `+2h` or `+1h30m`.

Todos are saved as you go in `~/.local/share/losenotion/todos.json` (or `$XDG_DATA_HOME/losenotion`; set `LOSENOTION_DATA_DIR` to use another folder).

## Reminders

Reminders are sent as desktop notifications by a separate background process, so they work when the app is closed:

```bash
losenotion reminders    # checks every 30s; -interval 10s to change
```

To start it automatically at login with systemd:

```bash
mkdir -p ~/.config/systemd/user
cp contrib/systemd/losenotion-reminders.service ~/.config/systemd/user/
systemctl --user enable --now losenotion-reminders
```

Reminders that come due while it isn't running are sent the next time it starts.
