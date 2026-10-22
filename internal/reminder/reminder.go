// Package reminder finds reminders that have come due and sends a
// notification for each. It runs as its own process (`losenotion reminders`)
// so reminders work without the TUI open.
//
// It only reads the todo list and keeps its own "last checked" time, so it
// never writes to the file the TUI is editing.
package reminder

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/hridayK/LoseNotion/internal/todo"
)

const notificationTitle = "LoseNotion reminder"

// Store is the part of store.Store the checker needs. Defining it here, where
// it is used, lets tests pass a fake.
type Store interface {
	Load() (todo.List, error)
	LastCheck() (time.Time, error)
	SetLastCheck(time.Time) error
}

type Notifier interface {
	Notify(title, body string) error
}

// Due returns the text of each todo that has a reminder in (after, upTo].
func Due(list todo.List, after, upTo time.Time) []string {
	var texts []string
	// reminders.REMIND.5: deleted todos live in list.Deleted, so they are never checked
	for _, t := range list.Todos {
		// reminders.REMIND.5
		if t.Done {
			continue
		}
		for _, at := range t.Reminders {
			if at.After(after) && !at.After(upTo) {
				texts = append(texts, t.Text)
				break // one notification per todo, even if several reminders are due
			}
		}
	}
	return texts
}

type Checker struct {
	Store    Store
	Notifier Notifier
	Now      func() time.Time
}

// Check notifies for every reminder that came due since the last check.
// reminders.REMIND.3
// reminders.ENG.3
func (c Checker) Check() error {
	now := c.Now()
	last, err := c.Store.LastCheck()
	if err != nil {
		return err
	}
	// The first ever check only records the time, so reminders from before
	// the process was set up don't all fire at once.
	if last.IsZero() || !now.After(last) {
		return c.Store.SetLastCheck(now)
	}

	list, err := c.Store.Load()
	if err != nil {
		return err
	}
	var errs []error
	for _, text := range Due(list, last, now) {
		if err := c.Notifier.Notify(notificationTitle, text); err != nil {
			errs = append(errs, err)
		}
	}
	// Move on even if a notification failed, so it isn't retried every tick.
	errs = append(errs, c.Store.SetLastCheck(now))
	return errors.Join(errs...)
}

// Run checks once immediately and then every interval until ctx is cancelled.
// Errors are logged and don't stop the loop.
// reminders.REMIND.4
// reminders.ENG.2
func (c Checker) Run(ctx context.Context, interval time.Duration, log io.Writer) error {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		if err := c.Check(); err != nil {
			fmt.Fprintln(log, "reminders:", err)
		}
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
		}
	}
}
