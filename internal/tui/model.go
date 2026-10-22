// Package tui is the Bubble Tea terminal UI for the todo list.
package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/hridayK/LoseNotion/internal/todo"
)

// mode is which screen is showing, and so which keys do what.
type mode int

const (
	modeList        mode = iota
	modeAdd              // typing a new todo
	modeDue              // typing the selected todo's due date
	modeTrash            // viewing deleted todos
	modeReminders        // viewing the selected todo's reminders
	modeAddReminder      // typing a new reminder time
)

var help = map[mode]string{
	modeList:        "a add · space done · d delete · u due date · r reminders · t deleted · q quit",
	modeAdd:         "enter add · esc cancel",
	modeDue:         "enter save · empty clears · esc cancel",
	modeTrash:       "r restore · x delete forever · esc back",
	modeReminders:   "a add · d remove · esc back",
	modeAddReminder: "enter add · esc cancel",
}

// todo-list.ENG.1
type Model struct {
	list todo.List
	save func(todo.List) error
	now  func() time.Time

	mode        mode
	cursor      int // selected todo
	trashCursor int // selected deleted todo
	remCursor   int // selected reminder of the selected todo
	input       textinput.Model
	status      string // last error, shown above the help line
}

// New returns a model showing list. save is called after every change, and
// now is the clock used for times like "+30m".
func New(list todo.List, save func(todo.List) error, now func() time.Time) Model {
	return Model{list: list, save: save, now: now, input: textinput.New()}
}

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// todo-list.QUIT.2
	if key, ok := msg.(tea.KeyMsg); ok && key.String() == "ctrl+c" {
		return m, tea.Quit
	}
	switch m.mode {
	case modeAdd, modeDue, modeAddReminder:
		return m.updatePrompt(msg)
	case modeTrash:
		return m.updateTrash(msg)
	case modeReminders:
		return m.updateReminders(msg)
	default:
		return m.updateList(msg)
	}
}

func (m Model) View() string {
	var b strings.Builder
	switch m.mode {
	case modeTrash:
		m.viewTrash(&b)
	case modeReminders, modeAddReminder:
		m.viewReminders(&b)
	default:
		m.viewList(&b)
	}
	if m.isPrompting() {
		fmt.Fprintf(&b, "\n%s\n", m.input.View())
	}
	if m.status != "" {
		fmt.Fprintf(&b, "\n%s\n", m.status)
	}
	fmt.Fprintf(&b, "\n%s\n", help[m.mode])
	return b.String()
}

// persist saves the list after a change. A failed save is shown to the user
// instead of being lost.
// persistence.SAVE.2
func (m Model) persist() Model {
	if m.save == nil {
		return m
	}
	if err := m.save(m.list); err != nil {
		m.status = "Couldn't save: " + err.Error()
	}
	return m
}

// moveCursor handles the up/down keys for a list of n items. moved is false
// if key isn't a movement key.
// todo-list.LIST.3
func moveCursor(key string, cursor, n int) (newCursor int, moved bool) {
	switch key {
	case "up", "k":
		return clamp(cursor-1, n), true
	case "down", "j":
		return clamp(cursor+1, n), true
	}
	return cursor, false
}

// clamp keeps a cursor inside a list of n items.
func clamp(cursor, n int) int {
	return max(0, min(cursor, n-1))
}

func cursorMark(selected bool) string {
	if selected {
		return "> "
	}
	return "  "
}
