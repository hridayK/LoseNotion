package tui

import (
	"testing"
	"time"
)

func TestReminders(t *testing.T) {
	t.Run("reminders.REMIND.1 adds several reminders to a todo", func(t *testing.T) {
		m := press(sampleModel(), "r",
			"a", "2026-10-23 08:00", "enter",
			"a", "+30m", "enter",
		)
		rs := m.list.Todos[0].Reminders
		want := []time.Time{testNow.Add(30 * time.Minute), time.Date(2026, 10, 23, 8, 0, 0, 0, time.Local)}
		if len(rs) != 2 || !rs[0].Equal(want[0]) || !rs[1].Equal(want[1]) {
			t.Fatalf("reminders = %v, want %v", rs, want)
		}
		if m.mode != modeReminders {
			t.Errorf("mode = %v, want back on the reminders screen", m.mode)
		}
		assertViewContains(t, m, `Reminders for "buy milk"`, "Thu 22 Oct 10:30", "Fri 23 Oct 08:00")
		if m = press(m, "esc"); m.mode != modeList {
			t.Errorf("mode after esc = %v, want the list", m.mode)
		}
		assertViewContains(t, m, "buy milk · 2 reminders")
	})

	t.Run("reminders.REMIND.1 rejects an unreadable time", func(t *testing.T) {
		m := press(sampleModel(), "r", "a", "soon", "enter")
		if len(m.list.Todos[0].Reminders) != 0 || m.mode != modeAddReminder {
			t.Errorf("reminders = %v, mode = %v; want none and the prompt still open", m.list.Todos[0].Reminders, m.mode)
		}
	})

	t.Run("reminders.REMIND.2 removes the selected reminder", func(t *testing.T) {
		m := press(sampleModel(), "r",
			"a", "+1h", "enter",
			"a", "+2h", "enter",
			"down", "d",
		)
		rs := m.list.Todos[0].Reminders
		if len(rs) != 1 || !rs[0].Equal(testNow.Add(time.Hour)) {
			t.Errorf("reminders = %v, want only the +1h one", rs)
		}
		if m.remCursor != 0 {
			t.Errorf("remCursor = %d, want 0", m.remCursor)
		}
	})
}
