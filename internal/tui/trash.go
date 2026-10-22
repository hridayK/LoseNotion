package tui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) updateTrash(msg tea.Msg) (tea.Model, tea.Cmd) {
	key, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}
	if c, moved := moveCursor(key.String(), m.trashCursor, len(m.list.Deleted)); moved {
		m.trashCursor = c
		return m, nil
	}

	if key.String() == "esc" || key.String() == "t" {
		m.mode = modeList
		return m, nil
	}
	if len(m.list.Deleted) == 0 {
		return m, nil
	}
	switch key.String() {
	// todo-list.DELETE.2-2
	case "r":
		m.list.Restore(m.trashCursor)
	// todo-list.DELETE.2-3
	case "x":
		m.list.Purge(m.trashCursor)
	default:
		return m, nil
	}
	m.trashCursor = clamp(m.trashCursor, len(m.list.Deleted))
	return m.persist(), nil
}

// todo-list.DELETE.2-1
func (m Model) viewTrash(b *strings.Builder) {
	b.WriteString("Deleted todos\n\n")
	if len(m.list.Deleted) == 0 {
		b.WriteString("Nothing here.\n")
		return
	}
	for i, t := range m.list.Deleted {
		b.WriteString(cursorMark(i == m.trashCursor) + t.Text + "\n")
	}
}
