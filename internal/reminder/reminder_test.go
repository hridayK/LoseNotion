package reminder

import (
	"context"
	"io"
	"slices"
	"testing"
	"time"

	"github.com/hridayK/LoseNotion/internal/store"
	"github.com/hridayK/LoseNotion/internal/todo"
)

type fakeNotifier struct{ bodies []string }

func (f *fakeNotifier) Notify(title, body string) error {
	f.bodies = append(f.bodies, body)
	return nil
}

// clock is a time source the test can move forward.
type clock struct{ now time.Time }

func (c *clock) Now() time.Time          { return c.now }
func (c *clock) Advance(d time.Duration) { c.now = c.now.Add(d) }

var start = time.Date(2026, 10, 23, 9, 0, 0, 0, time.UTC)

// setup saves list to a real store in a temp dir and returns a Checker that
// has already done its first check at start.
func setup(t *testing.T, list todo.List) (Checker, *fakeNotifier, *clock) {
	t.Helper()
	s, err := store.Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if err := s.Save(list); err != nil {
		t.Fatal(err)
	}
	n := &fakeNotifier{}
	clk := &clock{now: start}
	c := Checker{Store: s, Notifier: n, Now: clk.Now}
	if err := c.Check(); err != nil {
		t.Fatal(err)
	}
	return c, n, clk
}

func check(t *testing.T, c Checker) {
	t.Helper()
	if err := c.Check(); err != nil {
		t.Fatal(err)
	}
}

func withReminder(text string, at time.Time) todo.Todo {
	return todo.Todo{Text: text, Reminders: []time.Time{at}}
}

func TestRemind(t *testing.T) {
	t.Run("reminders.REMIND.3 notifies with the todo text when a reminder comes due", func(t *testing.T) {
		c, n, clk := setup(t, todo.List{Todos: []todo.Todo{
			withReminder("buy milk", start.Add(10*time.Minute)),
		}})

		clk.Advance(5 * time.Minute)
		check(t, c)
		if len(n.bodies) != 0 {
			t.Fatalf("notified %v before the reminder time", n.bodies)
		}

		clk.Advance(5 * time.Minute)
		check(t, c)
		if !slices.Equal(n.bodies, []string{"buy milk"}) {
			t.Errorf("notifications = %v, want [buy milk]", n.bodies)
		}
	})

	t.Run("reminders.REMIND.4 fires from saved data with no TUI running", func(t *testing.T) {
		// setup only touches the store and the checker: no tui.Model exists.
		c, n, clk := setup(t, todo.List{Todos: []todo.Todo{
			withReminder("stretch", start.Add(time.Minute)),
		}})
		clk.Advance(time.Minute)
		check(t, c)
		if len(n.bodies) != 1 {
			t.Errorf("notifications = %v, want 1", n.bodies)
		}
	})

	t.Run("reminders.REMIND.5 skips done and deleted todos", func(t *testing.T) {
		at := start.Add(time.Minute)
		done := withReminder("done one", at)
		done.Done = true
		c, n, clk := setup(t, todo.List{
			Todos:   []todo.Todo{done, withReminder("open one", at)},
			Deleted: []todo.Todo{withReminder("deleted one", at)},
		})
		clk.Advance(time.Minute)
		check(t, c)
		if !slices.Equal(n.bodies, []string{"open one"}) {
			t.Errorf("notifications = %v, want only [open one]", n.bodies)
		}
	})
}

func TestEng(t *testing.T) {
	t.Run("reminders.ENG.3 each reminder fires once", func(t *testing.T) {
		c, n, clk := setup(t, todo.List{Todos: []todo.Todo{
			withReminder("buy milk", start.Add(time.Minute)),
		}})
		for range 3 {
			clk.Advance(time.Minute)
			check(t, c)
		}
		if len(n.bodies) != 1 {
			t.Errorf("notifications = %v, want exactly 1", n.bodies)
		}
	})

	t.Run("reminders.ENG.3 reminders missed while stopped fire on the next start", func(t *testing.T) {
		c, n, clk := setup(t, todo.List{Todos: []todo.Todo{
			withReminder("missed", start.Add(time.Hour)),
		}})
		// Simulate the process being stopped for two hours: no checks happen.
		clk.Advance(2 * time.Hour)
		check(t, c)
		if !slices.Equal(n.bodies, []string{"missed"}) {
			t.Errorf("notifications = %v, want [missed]", n.bodies)
		}
	})

	t.Run("reminders.ENG.3 the first check does not fire old reminders", func(t *testing.T) {
		s, _ := store.Open(t.TempDir())
		s.Save(todo.List{Todos: []todo.Todo{withReminder("old", start.Add(-time.Hour))}})
		n := &fakeNotifier{}
		c := Checker{Store: s, Notifier: n, Now: (&clock{now: start}).Now}
		check(t, c)
		if len(n.bodies) != 0 {
			t.Errorf("first check notified %v, want nothing", n.bodies)
		}
	})

	t.Run("reminders.ENG.2 Run checks until cancelled", func(t *testing.T) {
		c, n, clk := setup(t, todo.List{Todos: []todo.Todo{
			withReminder("buy milk", start.Add(time.Minute)),
		}})
		clk.Advance(time.Minute)
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Run still does its immediate first check, then stops
		if err := c.Run(ctx, time.Hour, io.Discard); err != nil {
			t.Fatal(err)
		}
		if len(n.bodies) != 1 {
			t.Errorf("notifications = %v, want 1", n.bodies)
		}
	})
}

func TestDueWindowIsExclusiveThenInclusive(t *testing.T) {
	list := todo.List{Todos: []todo.Todo{
		withReminder("at after", start),
		withReminder("at upTo", start.Add(time.Minute)),
	}}
	got := Due(list, start, start.Add(time.Minute))
	if !slices.Equal(got, []string{"at upTo"}) {
		t.Errorf("Due = %v, want only the reminder at upTo", got)
	}
}
