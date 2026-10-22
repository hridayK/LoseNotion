package store

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/hridayK/LoseNotion/internal/todo"
)

func openTemp(t *testing.T) *Store {
	t.Helper()
	s, err := Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	return s
}

func TestSave(t *testing.T) {
	t.Run("persistence.SAVE.1 restores todos and deleted todos", func(t *testing.T) {
		s := openTemp(t)
		due := time.Date(2026, 10, 23, 18, 0, 0, 0, time.Local)
		want := todo.List{
			Todos: []todo.Todo{
				{Text: "buy milk", Due: due, Reminders: []time.Time{due.Add(-time.Hour)}},
				{Text: "call mum", Done: true},
			},
			Deleted: []todo.Todo{{Text: "old thing"}},
		}
		if err := s.Save(want); err != nil {
			t.Fatal(err)
		}

		reopened, err := Open(s.Dir())
		if err != nil {
			t.Fatal(err)
		}
		got, err := reopened.Load()
		if err != nil {
			t.Fatal(err)
		}
		if len(got.Todos) != 2 || len(got.Deleted) != 1 {
			t.Fatalf("got %+v, want 2 todos and 1 deleted", got)
		}
		first := got.Todos[0]
		if first.Text != "buy milk" || !first.Due.Equal(due) || len(first.Reminders) != 1 {
			t.Errorf("first todo = %+v, want text, due and reminder restored", first)
		}
		if !got.Todos[1].Done || got.Deleted[0].Text != "old thing" {
			t.Errorf("done flag or deleted todo not restored: %+v", got)
		}
	})

	t.Run("reminders.ENG.1 due dates and reminders survive a reload", func(t *testing.T) {
		s := openTemp(t)
		due := time.Date(2026, 10, 23, 18, 0, 0, 0, time.Local)
		reminders := []time.Time{due.Add(-2 * time.Hour), due.Add(-time.Hour)}
		if err := s.Save(todo.List{Todos: []todo.Todo{{Text: "a", Due: due, Reminders: reminders}}}); err != nil {
			t.Fatal(err)
		}
		got, err := s.Load()
		if err != nil {
			t.Fatal(err)
		}
		saved := got.Todos[0]
		if !saved.Due.Equal(due) || len(saved.Reminders) != 2 ||
			!saved.Reminders[0].Equal(reminders[0]) || !saved.Reminders[1].Equal(reminders[1]) {
			t.Errorf("reloaded %+v, want due %v and reminders %v", saved, due, reminders)
		}
	})

	t.Run("persistence.SAVE.1 a missing file is an empty list", func(t *testing.T) {
		got, err := openTemp(t).Load()
		if err != nil || len(got.Todos) != 0 {
			t.Errorf("Load() = %+v, %v; want empty list, nil", got, err)
		}
	})

	t.Run("persistence.SAVE.3 unreadable data is an error and is left alone", func(t *testing.T) {
		s := openTemp(t)
		path := filepath.Join(s.Dir(), todosFile)
		corrupt := []byte("{not json")
		if err := os.WriteFile(path, corrupt, 0o600); err != nil {
			t.Fatal(err)
		}
		if _, err := s.Load(); err == nil {
			t.Fatal("Load() succeeded on corrupt data, want error")
		}
		if data, _ := os.ReadFile(path); string(data) != string(corrupt) {
			t.Errorf("corrupt file was changed to %q", data)
		}
	})
}

func TestEng(t *testing.T) {
	t.Run("persistence.ENG.1 picks the data directory from the environment", func(t *testing.T) {
		t.Setenv("HOME", "/home/test")
		cases := []struct {
			override, xdg, want string
		}{
			{"/custom", "/xdg", "/custom"},
			{"", "/xdg", "/xdg/losenotion"},
			{"", "", "/home/test/.local/share/losenotion"},
		}
		for _, c := range cases {
			t.Setenv("LOSENOTION_DATA_DIR", c.override)
			t.Setenv("XDG_DATA_HOME", c.xdg)
			if got, err := DefaultDir(); err != nil || got != c.want {
				t.Errorf("DefaultDir() with override=%q xdg=%q = %q, %v; want %q", c.override, c.xdg, got, err, c.want)
			}
		}
	})

	t.Run("persistence.ENG.2 saves replace the file and leave no temp files", func(t *testing.T) {
		s := openTemp(t)
		for _, text := range []string{"first", "second"} {
			if err := s.Save(todo.List{Todos: []todo.Todo{{Text: text}}}); err != nil {
				t.Fatal(err)
			}
		}
		entries, _ := os.ReadDir(s.Dir())
		if len(entries) != 1 || entries[0].Name() != todosFile {
			var names []string
			for _, e := range entries {
				names = append(names, e.Name())
			}
			t.Errorf("data dir has %v, want only %s", names, todosFile)
		}
		if got, _ := s.Load(); got.Todos[0].Text != "second" {
			t.Errorf("Load() after two saves = %+v, want the second", got)
		}
	})
}

func TestLastCheckRoundTrip(t *testing.T) {
	s := openTemp(t)
	if got, err := s.LastCheck(); err != nil || !got.IsZero() {
		t.Fatalf("LastCheck() before any check = %v, %v; want zero time", got, err)
	}
	now := time.Now()
	if err := s.SetLastCheck(now); err != nil {
		t.Fatal(err)
	}
	if got, err := s.LastCheck(); err != nil || !got.Equal(now) {
		t.Errorf("LastCheck() = %v, %v; want %v", got, err, now)
	}
}
