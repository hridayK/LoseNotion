package tui

import "testing"

func TestAdd(t *testing.T) {
	t.Run("todo-list.ADD.1 adds a typed todo as not done", func(t *testing.T) {
		m, _ := newModel()
		m = press(m, "a", "buy milk", "enter")
		if len(m.list.Todos) != 1 {
			t.Fatalf("got %d todos, want 1", len(m.list.Todos))
		}
		if got := m.list.Todos[0]; got.Text != "buy milk" || got.Done {
			t.Errorf("got %+v, want {Text:buy milk Done:false}", got)
		}
		if m.mode != modeList {
			t.Errorf("mode = %v, want back in the list", m.mode)
		}
	})

	t.Run("todo-list.ADD.1 typed keys go into the text, not the list", func(t *testing.T) {
		m := press(sampleModel(), "a", "jkqd t", "enter")
		last := m.list.Todos[len(m.list.Todos)-1]
		if last.Text != "jkqd t" {
			t.Errorf("new todo text = %q, want %q", last.Text, "jkqd t")
		}
		if m.cursor != 0 || len(m.list.Deleted) != 0 {
			t.Errorf("list keys acted while typing: cursor = %d, deleted = %v", m.cursor, m.list.Deleted)
		}
	})

	t.Run("todo-list.ADD.2 cancelling does not create a todo", func(t *testing.T) {
		m, _ := newModel()
		m = press(m, "a", "buy milk", "esc")
		if len(m.list.Todos) != 0 {
			t.Errorf("got %d todos, want 0", len(m.list.Todos))
		}
		if m.mode != modeList {
			t.Errorf("mode = %v, want back in the list", m.mode)
		}
		if m = press(m, "a"); m.input.Value() != "" {
			t.Errorf("reopened prompt has leftover text %q", m.input.Value())
		}
	})

	t.Run("todo-list.ADD.3 does not add empty or whitespace-only text", func(t *testing.T) {
		for _, text := range []string{"", "   "} {
			m, _ := newModel()
			m = press(m, "a", text, "enter")
			if len(m.list.Todos) != 0 {
				t.Errorf("text %q: got %d todos, want 0", text, len(m.list.Todos))
			}
		}
	})
}
