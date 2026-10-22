package tui

import (
	"errors"
	"strings"
	"time"
)

// inputLayout is how dates are typed. Go layouts are written using the
// reference time Mon Jan 2 15:04:05 2006, not placeholders like YYYY.
const inputLayout = "2006-01-02 15:04"

var errBadTime = errors.New(`enter a time like "2026-10-23 18:00", or "+" and a duration like +30m or +2h`)

// parseWhen reads a time the user typed: "YYYY-MM-DD HH:MM" in local time,
// or "+" and a duration from now.
func parseWhen(s string, now time.Time) (time.Time, error) {
	s = strings.TrimSpace(s)
	if rest, ok := strings.CutPrefix(s, "+"); ok {
		d, err := time.ParseDuration(rest)
		if err != nil || d <= 0 {
			return time.Time{}, errBadTime
		}
		return now.Add(d).Truncate(time.Second), nil
	}
	t, err := time.ParseInLocation(inputLayout, s, now.Location())
	if err != nil {
		return time.Time{}, errBadTime
	}
	return t, nil
}

// formatInput formats t so parseWhen can read it back, for pre-filling a
// prompt. The zero time is an empty string.
func formatInput(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Local().Format(inputLayout)
}

// formatShow formats t for display, leaving out the year if it's this year.
func formatShow(t, now time.Time) string {
	t = t.Local()
	if t.Year() == now.Year() {
		return t.Format("Mon 2 Jan 15:04")
	}
	return t.Format("Mon 2 Jan 2006 15:04")
}
