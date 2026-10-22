// Package todo holds the todo list and the operations on it, with no
// knowledge of the terminal UI or of how the list is stored.
package todo

import (
	"errors"
	"slices"
	"strings"
	"time"
)

var ErrEmptyText = errors.New("todo text is empty")

// The json tags name each field in the saved file.
type Todo struct {
	Text string `json:"text"`
	Done bool   `json:"done"`
	// reminders.ENG.1
	Due       time.Time   `json:"due,omitzero"`
	Reminders []time.Time `json:"reminders,omitempty"`
}

// List is the user's todos plus the ones they deleted, which can be restored.
type List struct {
	Todos   []Todo `json:"todos"`
	Deleted []Todo `json:"deleted"`
}

// todo-list.ADD.1
func (l *List) Add(text string) error {
	text = strings.TrimSpace(text)
	// todo-list.ADD.3
	if text == "" {
		return ErrEmptyText
	}
	l.Todos = append(l.Todos, Todo{Text: text})
	return nil
}

// todo-list.COMPLETE.1
// todo-list.COMPLETE.2
func (l *List) ToggleDone(i int) {
	if inRange(l.Todos, i) {
		l.Todos[i].Done = !l.Todos[i].Done
	}
}

// todo-list.DELETE.1
// todo-list.DELETE.2
func (l *List) Delete(i int) {
	if inRange(l.Todos, i) {
		l.Deleted = append(l.Deleted, l.Todos[i])
		l.Todos = slices.Delete(l.Todos, i, i+1)
	}
}

// todo-list.DELETE.2-2
func (l *List) Restore(i int) {
	if inRange(l.Deleted, i) {
		l.Todos = append(l.Todos, l.Deleted[i])
		l.Deleted = slices.Delete(l.Deleted, i, i+1)
	}
}

// todo-list.DELETE.2-3
func (l *List) Purge(i int) {
	if inRange(l.Deleted, i) {
		l.Deleted = slices.Delete(l.Deleted, i, i+1)
	}
}

// SetDue sets the due date of todo i. The zero time clears it.
// reminders.DUE.1
// reminders.DUE.2
func (l *List) SetDue(i int, due time.Time) {
	if inRange(l.Todos, i) {
		l.Todos[i].Due = due
	}
}

// AddReminder adds a reminder to todo i, keeping reminders in time order.
// reminders.REMIND.1
func (l *List) AddReminder(i int, at time.Time) {
	if !inRange(l.Todos, i) {
		return
	}
	rs := append(l.Todos[i].Reminders, at)
	slices.SortFunc(rs, time.Time.Compare)
	l.Todos[i].Reminders = rs
}

// reminders.REMIND.2
func (l *List) RemoveReminder(i, r int) {
	if inRange(l.Todos, i) && inRange(l.Todos[i].Reminders, r) {
		l.Todos[i].Reminders = slices.Delete(l.Todos[i].Reminders, r, r+1)
	}
}

func inRange[T any](s []T, i int) bool {
	return i >= 0 && i < len(s)
}
