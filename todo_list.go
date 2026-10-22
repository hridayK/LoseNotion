package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type todo struct {
	text string
	done bool
}

// mode is which view the todo list is showing, and so which keys it accepts.
type mode int

const (
	modeList mode = iota
	modeAdding
)

// todo-list.ENG.1
// todo-list.ENG.2
type todoList struct {
	todos  []todo
	cursor int // index of the selected todo
	mode   mode
	input  textinput.Model // text of the todo being added
}

func newTodoList(todos ...todo) todoList {
	input := textinput.New()
	input.Prompt = "New todo: "
	return todoList{todos: todos, input: input}
}

func (m todoList) Init() tea.Cmd { return nil }

func (m todoList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// todo-list.QUIT.2
	if key, ok := msg.(tea.KeyMsg); ok && key.String() == "ctrl+c" {
		return m, tea.Quit
	}
	if m.mode == modeAdding {
		return m.updateAdding(msg)
	}
	return m.updateList(msg)
}

func (m todoList) updateList(msg tea.Msg) (tea.Model, tea.Cmd) {
	key, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}
	switch key.String() {
	// todo-list.LIST.3
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}
	case "down", "j":
		if m.cursor < len(m.todos)-1 {
			m.cursor++
		}
	// todo-list.ADD.1
	case "a":
		m.mode = modeAdding
		return m, m.input.Focus()
	}
	return m, nil
}

func (m todoList) updateAdding(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok {
		switch key.String() {
		// todo-list.ADD.1
		case "enter":
			text := strings.TrimSpace(m.input.Value())
			// todo-list.ADD.3
			if text == "" {
				return m, nil
			}
			m.todos = append(m.todos, todo{text: text})
			return m.stopAdding(), nil
		// todo-list.ADD.2
		case "esc":
			return m.stopAdding(), nil
		}
	}

	// Everything else (typed keys, cursor blinks) belongs to the text input.
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m todoList) stopAdding() todoList {
	m.input.Reset()
	m.input.Blur()
	m.mode = modeList
	return m
}

func (m todoList) View() string {
	// todo-list.LIST.2
	if len(m.todos) == 0 && m.mode == modeList {
		return "No todos yet. Press a to add one.\n"
	}

	// todo-list.LIST.1
	var b strings.Builder
	for i, t := range m.todos {
		cursor := "  "
		if i == m.cursor {
			cursor = "> "
		}
		// todo-list.LIST.4
		check := "[ ]"
		if t.done {
			check = "[x]"
		}
		fmt.Fprintf(&b, "%s%s %s\n", cursor, check, t.text)
	}

	if m.mode == modeAdding {
		fmt.Fprintf(&b, "\n%s\n", m.input.View())
	}
	return b.String()
}
