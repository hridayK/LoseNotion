package tui

import (
	"slices"
	"testing"
)

func TestList(t *testing.T) {
	t.Run("todo-list.ENG.1 Update works without a terminal and returns a new model", func(t *testing.T) {
		before := sampleModel()
		after := press(before, "down")
		if before.cursor != 0 || after.cursor != 1 {
			t.Errorf("cursor before = %d, after = %d; want 0 then 1", before.cursor, after.cursor)
		}
	})

	t.Run("todo-list.LIST.1 shows the todos on launch", func(t *testing.T) {
		assertViewContains(t, sampleModel(), "buy milk", "call mum", "write spec")
	})

	t.Run("todo-list.LIST.2 shows an empty-state message", func(t *testing.T) {
		m, _ := newModel()
		assertViewContains(t, m, "No todos yet")
	})

	t.Run("todo-list.LIST.3 moves the selection down and up", func(t *testing.T) {
		m := sampleModel()
		steps := []struct {
			key  string
			want int
		}{
			{"down", 1}, {"j", 2}, {"up", 1}, {"k", 0},
		}
		for _, s := range steps {
			m = press(m, s.key)
			if m.cursor != s.want {
				t.Fatalf("after %q: cursor = %d, want %d", s.key, m.cursor, s.want)
			}
		}
	})

	t.Run("todo-list.LIST.3 keeps the selection inside the list", func(t *testing.T) {
		m := press(sampleModel(), "up")
		if m.cursor != 0 {
			t.Errorf("up at top: cursor = %d, want 0", m.cursor)
		}
		m = press(m, "down", "down", "down", "down")
		if m.cursor != 2 {
			t.Errorf("down past bottom: cursor = %d, want 2", m.cursor)
		}
	})

	t.Run("todo-list.LIST.4 shows whether each todo is done", func(t *testing.T) {
		assertViewContains(t, sampleModel(), "[ ] buy milk", "[x] call mum")
	})
}

func TestComplete(t *testing.T) {
	t.Run("todo-list.COMPLETE.1 marks the selected todo as done", func(t *testing.T) {
		m := press(sampleModel(), "space")
		if !m.list.Todos[0].Done {
			t.Error("selected todo is not done after pressing space")
		}
		assertViewContains(t, m, "[x] buy milk")
	})

	t.Run("todo-list.COMPLETE.2 marks a done todo as not done again", func(t *testing.T) {
		m := press(sampleModel(), "down", "space") // "call mum" starts done
		if m.list.Todos[1].Done {
			t.Error("done todo is still done after pressing space")
		}
	})
}

func TestDelete(t *testing.T) {
	t.Run("todo-list.DELETE.1 deletes the selected todo", func(t *testing.T) {
		m := press(sampleModel(), "down", "d")
		if got := todoTexts(m.list.Todos); !slices.Equal(got, []string{"buy milk", "write spec"}) {
			t.Errorf("todos = %v, want call mum removed", got)
		}
	})

	t.Run("todo-list.DELETE.1 keeps the selection inside the list after deleting the last todo", func(t *testing.T) {
		m := press(sampleModel(), "down", "down", "d")
		if m.cursor != 1 {
			t.Errorf("cursor = %d, want 1", m.cursor)
		}
	})

	t.Run("todo-list.DELETE.2 moves the deleted todo to the deleted list", func(t *testing.T) {
		m := press(sampleModel(), "d")
		if got := todoTexts(m.list.Deleted); !slices.Equal(got, []string{"buy milk"}) {
			t.Errorf("deleted = %v, want [buy milk]", got)
		}
	})
}

func TestQuit(t *testing.T) {
	t.Run("todo-list.QUIT.1 q quits from the list view", func(t *testing.T) {
		if !quits(sampleModel(), "q") {
			t.Error("q did not quit")
		}
	})

	t.Run("todo-list.QUIT.1 q is typed, not quit, while adding a todo", func(t *testing.T) {
		m := press(sampleModel(), "a", "q")
		if m.input.Value() != "q" || m.mode != modeAdd {
			t.Errorf("input = %q, mode = %v; want q typed into the add prompt", m.input.Value(), m.mode)
		}
	})

	t.Run("todo-list.QUIT.2 ctrl+c quits from every view", func(t *testing.T) {
		views := map[string][]string{
			"list":         nil,
			"add":          {"a"},
			"due date":     {"u"},
			"deleted":      {"t"},
			"reminders":    {"r"},
			"add reminder": {"r", "a"},
		}
		for name, keys := range views {
			if !quits(press(sampleModel(), keys...), "ctrl+c") {
				t.Errorf("ctrl+c did not quit from the %s view", name)
			}
		}
	})
}
