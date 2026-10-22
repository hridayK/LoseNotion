package tui

import (
	"testing"
	"time"
)

func TestDue(t *testing.T) {
	t.Run("reminders.DUE.1 sets a due date and time", func(t *testing.T) {
		m := press(sampleModel(), "u", "2026-10-23 18:30", "enter")
		want := time.Date(2026, 10, 23, 18, 30, 0, 0, time.Local)
		if got := m.list.Todos[0].Due; !got.Equal(want) {
			t.Errorf("due = %v, want %v", got, want)
		}
		if m.mode != modeList {
			t.Errorf("mode = %v, want back in the list", m.mode)
		}
	})

	t.Run("reminders.DUE.1 accepts a time relative to now", func(t *testing.T) {
		m := press(sampleModel(), "u", "+2h", "enter")
		if got, want := m.list.Todos[0].Due, testNow.Add(2*time.Hour); !got.Equal(want) {
			t.Errorf("due = %v, want %v", got, want)
		}
	})

	t.Run("reminders.DUE.1 rejects an unreadable time and keeps the prompt open", func(t *testing.T) {
		m := press(sampleModel(), "u", "tomorrow-ish", "enter")
		if !m.list.Todos[0].Due.IsZero() || m.mode != modeDue {
			t.Errorf("due = %v, mode = %v; want no due date and the prompt still open", m.list.Todos[0].Due, m.mode)
		}
		assertViewContains(t, m, "enter a time like")
	})

	t.Run("reminders.DUE.2 changes a due date, starting from the current one", func(t *testing.T) {
		m := press(sampleModel(), "u", "2026-10-23 18:30", "enter", "u")
		if got := m.input.Value(); got != "2026-10-23 18:30" {
			t.Fatalf("prompt prefilled with %q, want the current due date", got)
		}
		m.input.SetValue("2026-10-24 09:00")
		m = press(m, "enter")
		want := time.Date(2026, 10, 24, 9, 0, 0, 0, time.Local)
		if got := m.list.Todos[0].Due; !got.Equal(want) {
			t.Errorf("due = %v, want %v", got, want)
		}
	})

	t.Run("reminders.DUE.2 clears a due date with an empty value", func(t *testing.T) {
		m := press(sampleModel(), "u", "+1h", "enter", "u")
		m.input.SetValue("")
		m = press(m, "enter")
		if got := m.list.Todos[0].Due; !got.IsZero() {
			t.Errorf("due = %v, want cleared", got)
		}
	})

	t.Run("reminders.DUE.3 shows the due date in the list", func(t *testing.T) {
		m := press(sampleModel(), "u", "2026-10-23 18:30", "enter")
		assertViewContains(t, m, "buy milk  (due Fri 23 Oct 18:30)")
	})
}
