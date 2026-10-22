package tui

import (
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/hridayK/LoseNotion/internal/todo"
)

// testNow is the fixed "current time" for every test model.
var testNow = time.Date(2026, 10, 22, 10, 0, 0, 0, time.Local)

// Keys that Bubble Tea reports as a KeyType rather than as typed runes.
var specialKeys = map[string]tea.KeyType{
	"up":     tea.KeyUp,
	"down":   tea.KeyDown,
	"enter":  tea.KeyEnter,
	"esc":    tea.KeyEsc,
	"space":  tea.KeySpace,
	"ctrl+c": tea.KeyCtrlC,
}

func keyMsg(key string) tea.KeyMsg {
	switch kt, ok := specialKeys[key]; {
	case kt == tea.KeySpace:
		return tea.KeyMsg{Type: kt, Runes: []rune{' '}}
	case ok:
		return tea.KeyMsg{Type: kt}
	}
	return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(key)}
}

// press sends keys to the model one at a time, the same way Bubble Tea would,
// and returns the updated model. Strings that aren't special keys are typed
// as text.
func press(m Model, keys ...string) Model {
	for _, k := range keys {
		next, _ := m.Update(keyMsg(k))
		m = next.(Model)
	}
	return m
}

// quits reports whether pressing key asks Bubble Tea to quit.
func quits(m Model, key string) bool {
	_, cmd := m.Update(keyMsg(key))
	if cmd == nil {
		return false
	}
	_, ok := cmd().(tea.QuitMsg)
	return ok
}

// saveLog records every list the model saves.
type saveLog struct{ lists []todo.List }

func (s *saveLog) save(l todo.List) error {
	s.lists = append(s.lists, l)
	return nil
}

func newModel(todos ...todo.Todo) (Model, *saveLog) {
	log := &saveLog{}
	m := New(todo.List{Todos: todos}, log.save, func() time.Time { return testNow })
	return m, log
}

func sampleModel() Model {
	m, _ := newModel(
		todo.Todo{Text: "buy milk"},
		todo.Todo{Text: "call mum", Done: true},
		todo.Todo{Text: "write spec"},
	)
	return m
}

func todoTexts(ts []todo.Todo) []string {
	var out []string
	for _, t := range ts {
		out = append(out, t.Text)
	}
	return out
}

func assertViewContains(t *testing.T, m Model, wants ...string) {
	t.Helper()
	view := m.View()
	for _, want := range wants {
		if !strings.Contains(view, want) {
			t.Errorf("view is missing %q:\n%s", want, view)
		}
	}
}
