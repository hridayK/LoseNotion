package tui

import (
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// The add-todo, due-date and add-reminder screens all share one text input.

func (m Model) isPrompting() bool {
	return m.mode == modeAdd || m.mode == modeDue || m.mode == modeAddReminder
}

func (m Model) openPrompt(md mode, label, value, placeholder string) (tea.Model, tea.Cmd) {
	m.mode = md
	m.status = ""
	m.input.Prompt = label
	m.input.Placeholder = placeholder
	m.input.SetValue(value)
	return m, m.input.Focus()
}

// closePrompt goes back to the screen the prompt was opened from.
func (m Model) closePrompt() Model {
	m.input.Reset()
	m.input.Blur()
	if m.mode == modeAddReminder {
		m.mode = modeReminders
	} else {
		m.mode = modeList
	}
	return m
}

func (m Model) updatePrompt(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok {
		switch key.String() {
		case "enter":
			return m.submitPrompt()
		// todo-list.ADD.2
		case "esc":
			m.status = ""
			return m.closePrompt(), nil
		}
	}

	// Everything else (typed keys, cursor blinks) belongs to the text input.
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

// submitPrompt applies what was typed. On invalid input the prompt stays
// open so the user can fix it.
func (m Model) submitPrompt() (tea.Model, tea.Cmd) {
	value := m.input.Value()
	switch m.mode {
	case modeAdd:
		// todo-list.ADD.1
		// todo-list.ADD.3
		if err := m.list.Add(value); err != nil {
			return m, nil
		}
	case modeDue:
		// reminders.DUE.2: an empty value clears the due date
		var due time.Time
		if strings.TrimSpace(value) != "" {
			t, err := parseWhen(value, m.now())
			if err != nil {
				m.status = err.Error()
				return m, nil
			}
			due = t
		}
		// reminders.DUE.1
		m.list.SetDue(m.cursor, due)
	case modeAddReminder:
		at, err := parseWhen(value, m.now())
		if err != nil {
			m.status = err.Error()
			return m, nil
		}
		// reminders.REMIND.1
		m.list.AddReminder(m.cursor, at)
	}
	m.status = ""
	return m.closePrompt().persist(), nil
}
