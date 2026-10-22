// Package store saves the todo list and reminder state as files in a data
// directory.
package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/hridayK/LoseNotion/internal/todo"
)

const (
	todosFile     = "todos.json"
	lastCheckFile = "reminders-last-check"
)

type Store struct {
	dir string
}

// DefaultDir is where data lives unless the user overrides it.
// persistence.ENG.1
func DefaultDir() (string, error) {
	if dir := os.Getenv("LOSENOTION_DATA_DIR"); dir != "" {
		return dir, nil
	}
	if xdg := os.Getenv("XDG_DATA_HOME"); xdg != "" {
		return filepath.Join(xdg, "losenotion"), nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".local", "share", "losenotion"), nil
}

// Open returns a Store for dir, creating the directory if needed.
func Open(dir string) (*Store, error) {
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return nil, err
	}
	return &Store{dir: dir}, nil
}

func (s *Store) Dir() string { return s.dir }

// Load reads the saved list. A missing file is an empty list, not an error.
// persistence.SAVE.1
func (s *Store) Load() (todo.List, error) {
	var list todo.List
	data, err := os.ReadFile(filepath.Join(s.dir, todosFile))
	if errors.Is(err, fs.ErrNotExist) {
		return list, nil
	}
	if err != nil {
		return list, err
	}
	// persistence.SAVE.3
	if err := json.Unmarshal(data, &list); err != nil {
		return list, fmt.Errorf("reading %s: %w", filepath.Join(s.dir, todosFile), err)
	}
	return list, nil
}

// persistence.SAVE.2
func (s *Store) Save(list todo.List) error {
	data, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		return err
	}
	return s.writeAtomic(todosFile, data)
}

// LastCheck is when the reminder process last checked for due reminders, or
// the zero time if it never has.
func (s *Store) LastCheck() (time.Time, error) {
	data, err := os.ReadFile(filepath.Join(s.dir, lastCheckFile))
	if errors.Is(err, fs.ErrNotExist) {
		return time.Time{}, nil
	}
	if err != nil {
		return time.Time{}, err
	}
	return time.Parse(time.RFC3339Nano, strings.TrimSpace(string(data)))
}

func (s *Store) SetLastCheck(t time.Time) error {
	return s.writeAtomic(lastCheckFile, []byte(t.Format(time.RFC3339Nano)+"\n"))
}

// writeAtomic writes to a temp file and renames it into place. A rename
// replaces the file in one step, so a reader (such as the reminder process)
// sees either the old file or the new one, never half of each.
// persistence.ENG.2
func (s *Store) writeAtomic(name string, data []byte) error {
	tmp, err := os.CreateTemp(s.dir, name+".*.tmp")
	if err != nil {
		return err
	}
	defer os.Remove(tmp.Name()) // no-op once the rename succeeds

	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	return os.Rename(tmp.Name(), filepath.Join(s.dir, name))
}
