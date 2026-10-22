package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// The reminders screen always belongs to the todo selected in the list
// (m.cursor), which can't change while the screen is open.

func (m Model) updateReminders(msg tea.Msg) (tea.Model, tea.Cmd) {
	key, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}
	m.status = ""
	reminders := m.list.Todos[m.cursor].Reminders
	if c, moved := moveCursor(key.String(), m.remCursor, len(reminders)); moved {
		m.remCursor = c
		return m, nil
	}

	switch key.String() {
	case "esc":
		m.mode = modeList
	// reminders.REMIND.1
	case "a":
		return m.openPrompt(modeAddReminder, "Remind at: ", "", "2026-10-23 18:00 or +30m")
	// reminders.REMIND.2
	case "d":
		if len(reminders) > 0 {
			m.list.RemoveReminder(m.cursor, m.remCursor)
			m.remCursor = clamp(m.remCursor, len(m.list.Todos[m.cursor].Reminders))
			return m.persist(), nil
		}
	}
	return m, nil
}

func (m Model) viewReminders(b *strings.Builder) {
	t := m.list.Todos[m.cursor]
	fmt.Fprintf(b, "Reminders for %q\n\n", t.Text)
	if len(t.Reminders) == 0 {
		b.WriteString("No reminders yet.\n")
		return
	}
	for i, at := range t.Reminders {
		b.WriteString(cursorMark(i == m.remCursor) + formatShow(at, m.now()) + "\n")
	}
}
