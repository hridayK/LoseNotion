package tui

import (
	"errors"
	"slices"
	"testing"
	"time"

	"github.com/hridayK/LoseNotion/internal/todo"
)

func TestPersist(t *testing.T) {
	t.Run("persistence.SAVE.2 every change is saved straight away", func(t *testing.T) {
		m, log := newModel(todo.Todo{Text: "buy milk"})
		changes := [][]string{
			{"a", "call mum", "enter"},
			{"space"},
			{"u", "+1h", "enter"},
			{"r", "a", "+30m", "enter", "esc"},
			{"d"},
			{"t", "r", "esc"},
		}
		for i, keys := range changes {
			m = press(m, keys...)
			if len(log.lists) != i+1 {
				t.Fatalf("after %v: %d saves, want %d", keys, len(log.lists), i+1)
			}
		}
		saved := log.lists[len(log.lists)-1]
		if !slices.Equal(todoTexts(saved.Todos), todoTexts(m.list.Todos)) {
			t.Errorf("last save = %v, want the current list %v", todoTexts(saved.Todos), todoTexts(m.list.Todos))
		}
	})

	t.Run("persistence.SAVE.2 moving around does not save", func(t *testing.T) {
		m, log := newModel(todo.Todo{Text: "a"}, todo.Todo{Text: "b"})
		press(m, "down", "up", "t", "esc", "a", "esc")
		if len(log.lists) != 0 {
			t.Errorf("got %d saves, want 0", len(log.lists))
		}
	})

	t.Run("persistence.SAVE.2 keys with nothing to act on do not save", func(t *testing.T) {
		empty, emptyLog := newModel()
		press(empty, "space", "d", "t", "r", "x", "esc")
		noReminders, remLog := newModel(todo.Todo{Text: "a"})
		press(noReminders, "r", "d", "esc")
		if n := len(emptyLog.lists) + len(remLog.lists); n != 0 {
			t.Errorf("got %d saves, want 0", n)
		}
	})

	t.Run("persistence.SAVE.2 a failed save is shown", func(t *testing.T) {
		failing := func(todo.List) error { return errors.New("disk full") }
		m := New(todo.List{}, failing, func() time.Time { return testNow })
		m = press(m, "a", "buy milk", "enter")
		assertViewContains(t, m, "Couldn't save: disk full")
	})
}
