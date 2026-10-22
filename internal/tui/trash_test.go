package tui

import (
	"slices"
	"testing"
)

// withDeleted returns the sample model with "buy milk" and "call mum" deleted,
// showing the deleted list.
func withDeleted() Model {
	return press(sampleModel(), "d", "d", "t")
}

func TestTrash(t *testing.T) {
	t.Run("todo-list.DELETE.2-1 shows the deleted todos", func(t *testing.T) {
		m := withDeleted()
		if m.mode != modeTrash {
			t.Fatalf("mode = %v, want the deleted list", m.mode)
		}
		assertViewContains(t, m, "Deleted todos", "buy milk", "call mum")
	})

	t.Run("todo-list.DELETE.2-1 shows an empty message and goes back with esc", func(t *testing.T) {
		m := press(sampleModel(), "t")
		assertViewContains(t, m, "Nothing here")
		if m = press(m, "esc"); m.mode != modeList {
			t.Errorf("mode after esc = %v, want the list", m.mode)
		}
	})

	t.Run("todo-list.DELETE.2-2 restores a deleted todo to the list", func(t *testing.T) {
		m := press(withDeleted(), "down", "r") // restore "call mum"
		if got := todoTexts(m.list.Todos); !slices.Equal(got, []string{"write spec", "call mum"}) {
			t.Errorf("todos = %v, want call mum restored", got)
		}
		if got := todoTexts(m.list.Deleted); !slices.Equal(got, []string{"buy milk"}) {
			t.Errorf("deleted = %v, want only buy milk left", got)
		}
		if !m.list.Todos[1].Done {
			t.Error("restored todo lost its done state")
		}
	})

	t.Run("todo-list.DELETE.2-3 permanently deletes a deleted todo", func(t *testing.T) {
		m := press(withDeleted(), "x")
		if got := todoTexts(m.list.Deleted); !slices.Equal(got, []string{"call mum"}) {
			t.Errorf("deleted = %v, want buy milk gone", got)
		}
		if got := todoTexts(m.list.Todos); !slices.Equal(got, []string{"write spec"}) {
			t.Errorf("todos = %v, want unchanged", got)
		}
	})
}
