package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) updateList(msg tea.Msg) (tea.Model, tea.Cmd) {
	key, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}
	m.status = ""
	if c, moved := moveCursor(key.String(), m.cursor, len(m.list.Todos)); moved {
		m.cursor = c
		return m, nil
	}

	hasTodo := len(m.list.Todos) > 0
	switch key.String() {
	// todo-list.QUIT.1
	case "q":
		return m, tea.Quit
	// todo-list.ADD.1
	case "a":
		return m.openPrompt(modeAdd, "New todo: ", "", "")
	// todo-list.DELETE.2-1
	case "t":
		m.mode = modeTrash
		m.trashCursor = 0
	// todo-list.COMPLETE.1
	// todo-list.COMPLETE.2
	case " ":
		if hasTodo {
			m.list.ToggleDone(m.cursor)
			return m.persist(), nil
		}
	// todo-list.DELETE.1
	case "d":
		if hasTodo {
			m.list.Delete(m.cursor)
			m.cursor = clamp(m.cursor, len(m.list.Todos))
			return m.persist(), nil
		}
	// reminders.DUE.1
	// reminders.DUE.2
	case "u":
		if hasTodo {
			due := formatInput(m.list.Todos[m.cursor].Due)
			return m.openPrompt(modeDue, "Due: ", due, "2026-10-23 18:00 or +2h")
		}
	// reminders.REMIND.1
	case "r":
		if hasTodo {
			m.mode = modeReminders
			m.remCursor = 0
		}
	}
	return m, nil
}

func (m Model) viewList(b *strings.Builder) {
	// todo-list.LIST.2
	if len(m.list.Todos) == 0 {
		b.WriteString("No todos yet.\n")
		return
	}

	// todo-list.LIST.1
	for i, t := range m.list.Todos {
		// todo-list.LIST.4
		check := "[ ]"
		if t.Done {
			check = "[x]"
		}
		fmt.Fprintf(b, "%s%s %s", cursorMark(i == m.cursor), check, t.Text)
		// reminders.DUE.3
		if !t.Due.IsZero() {
			fmt.Fprintf(b, "  (due %s)", formatShow(t.Due, m.now()))
		}
		if n := len(t.Reminders); n == 1 {
			b.WriteString(" · 1 reminder")
		} else if n > 1 {
			fmt.Fprintf(b, " · %d reminders", n)
		}
		b.WriteString("\n")
	}
}
