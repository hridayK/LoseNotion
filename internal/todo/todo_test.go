package todo

import (
	"errors"
	"testing"
	"time"
)

func texts(ts []Todo) []string {
	var out []string
	for _, t := range ts {
		out = append(out, t.Text)
	}
	return out
}

func newList(items ...string) List {
	var l List
	for _, s := range items {
		l.Add(s)
	}
	return l
}

func TestAddTrimsAndRejectsBlankText(t *testing.T) {
	l := newList("  buy milk  ")
	if got := l.Todos[0].Text; got != "buy milk" {
		t.Errorf("text = %q, want trimmed %q", got, "buy milk")
	}
	if err := l.Add(" \t "); !errors.Is(err, ErrEmptyText) {
		t.Errorf("Add(blank) error = %v, want ErrEmptyText", err)
	}
}

func TestDeleteRestorePurge(t *testing.T) {
	l := newList("a", "b", "c")
	l.Delete(1)
	if got := texts(l.Todos); len(got) != 2 || got[0] != "a" || got[1] != "c" {
		t.Fatalf("todos after delete = %v, want [a c]", got)
	}
	l.Delete(0)
	l.Restore(0) // "b", the first one deleted
	if got := texts(l.Todos); len(got) != 2 || got[1] != "b" {
		t.Fatalf("todos after restore = %v, want b appended", got)
	}
	l.Purge(0)
	if len(l.Deleted) != 0 {
		t.Errorf("deleted after purge = %v, want empty", texts(l.Deleted))
	}
}

func TestOutOfRangeIndexesAreIgnored(t *testing.T) {
	var l List
	l.ToggleDone(0)
	l.Delete(-1)
	l.Restore(3)
	l.Purge(0)
	l.SetDue(0, time.Now())
	l.AddReminder(0, time.Now())
	l.RemoveReminder(0, 0)
	if len(l.Todos) != 0 || len(l.Deleted) != 0 {
		t.Errorf("empty list changed: %+v", l)
	}
}

func TestAddReminderKeepsTimeOrder(t *testing.T) {
	l := newList("a")
	base := time.Date(2026, 10, 23, 9, 0, 0, 0, time.UTC)
	l.AddReminder(0, base.Add(2*time.Hour))
	l.AddReminder(0, base)
	l.AddReminder(0, base.Add(time.Hour))
	rs := l.Todos[0].Reminders
	for i := 1; i < len(rs); i++ {
		if rs[i].Before(rs[i-1]) {
			t.Fatalf("reminders not sorted: %v", rs)
		}
	}
}
