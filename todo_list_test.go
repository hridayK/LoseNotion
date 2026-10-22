package main

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

// Keys that Bubble Tea reports as a KeyType rather than as typed runes.
var specialKeys = map[string]tea.KeyType{
	"up":    tea.KeyUp,
	"down":  tea.KeyDown,
	"enter": tea.KeyEnter,
	"esc":   tea.KeyEsc,
}

// press sends one key to the model, the same way Bubble Tea would, and
// returns the updated model. Any other string is sent as typed text.
func press(m todoList, key string) todoList {
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(key)}
	if kt, ok := specialKeys[key]; ok {
		msg = tea.KeyMsg{Type: kt}
	}
	next, _ := m.Update(msg)
	return next.(todoList)
}

func sampleTodos() []todo {
	return []todo{
		{text: "buy milk"},
		{text: "call mum", done: true},
		{text: "write spec"},
	}
}

func TestList(t *testing.T) {
	t.Run("todo-list.LIST.1 shows the todos on launch", func(t *testing.T) {
		view := newTodoList(sampleTodos()...).View()
		for _, td := range sampleTodos() {
			if !strings.Contains(view, td.text) {
				t.Errorf("view is missing %q:\n%s", td.text, view)
			}
		}
	})

	t.Run("todo-list.LIST.2 shows an empty-state message", func(t *testing.T) {
		view := newTodoList().View()
		if !strings.Contains(view, "No todos yet") {
			t.Errorf("expected empty-state message, got:\n%s", view)
		}
	})

	t.Run("todo-list.LIST.3 moves the selection down and up", func(t *testing.T) {
		m := newTodoList(sampleTodos()...)
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
		m := press(newTodoList(sampleTodos()...), "up")
		if m.cursor != 0 {
			t.Errorf("up at top: cursor = %d, want 0", m.cursor)
		}
		for range 5 {
			m = press(m, "down")
		}
		if m.cursor != 2 {
			t.Errorf("down past bottom: cursor = %d, want 2", m.cursor)
		}
	})

	t.Run("todo-list.LIST.4 shows whether each todo is done", func(t *testing.T) {
		view := newTodoList(sampleTodos()...).View()
		for _, want := range []string{"[ ] buy milk", "[x] call mum"} {
			if !strings.Contains(view, want) {
				t.Errorf("view is missing %q:\n%s", want, view)
			}
		}
	})
}

// addTodo opens the add prompt, types text, and presses enter.
func addTodo(m todoList, text string) todoList {
	return press(press(press(m, "a"), text), "enter")
}

func TestAdd(t *testing.T) {
	t.Run("todo-list.ADD.1 adds a typed todo as not done", func(t *testing.T) {
		m := addTodo(newTodoList(), "buy milk")
		if len(m.todos) != 1 {
			t.Fatalf("got %d todos, want 1", len(m.todos))
		}
		if got := m.todos[0]; got.text != "buy milk" || got.done {
			t.Errorf("got %+v, want {text:buy milk done:false}", got)
		}
		if m.mode != modeList {
			t.Errorf("mode = %v, want back in the list", m.mode)
		}
	})

	t.Run("todo-list.ADD.1 typed keys go into the text, not the list", func(t *testing.T) {
		m := addTodo(newTodoList(sampleTodos()...), "jka")
		last := m.todos[len(m.todos)-1]
		if last.text != "jka" {
			t.Errorf("new todo text = %q, want %q", last.text, "jka")
		}
		if m.cursor != 0 {
			t.Errorf("cursor = %d, want 0 (j/k should not move it while typing)", m.cursor)
		}
	})

	t.Run("todo-list.ADD.2 cancelling does not create a todo", func(t *testing.T) {
		m := press(press(press(newTodoList(), "a"), "buy milk"), "esc")
		if len(m.todos) != 0 {
			t.Errorf("got %d todos, want 0", len(m.todos))
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
			m := press(press(press(newTodoList(), "a"), text), "enter")
			if len(m.todos) != 0 {
				t.Errorf("text %q: got %d todos, want 0", text, len(m.todos))
			}
		}
	})
}
