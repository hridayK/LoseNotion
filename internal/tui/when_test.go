package tui

import (
	"testing"
	"time"
)

func TestParseWhen(t *testing.T) {
	cases := []struct {
		in   string
		want time.Time
		ok   bool
	}{
		{"2026-10-23 18:30", time.Date(2026, 10, 23, 18, 30, 0, 0, time.Local), true},
		{"  2026-10-23 18:30  ", time.Date(2026, 10, 23, 18, 30, 0, 0, time.Local), true},
		{"+90m", testNow.Add(90 * time.Minute), true},
		{"+1h30m", testNow.Add(90 * time.Minute), true},
		{"+0s", time.Time{}, false},
		{"+-1h", time.Time{}, false},
		{"+1d", time.Time{}, false}, // Go durations have no days unit
		{"2026-10-23", time.Time{}, false},
		{"tomorrow", time.Time{}, false},
	}
	for _, c := range cases {
		got, err := parseWhen(c.in, testNow)
		if (err == nil) != c.ok || !got.Equal(c.want) {
			t.Errorf("parseWhen(%q) = %v, %v; want %v, ok=%v", c.in, got, err, c.want, c.ok)
		}
	}
}

func TestFormatInputRoundTrips(t *testing.T) {
	due := time.Date(2026, 10, 23, 18, 30, 0, 0, time.Local)
	got, err := parseWhen(formatInput(due), testNow)
	if err != nil || !got.Equal(due) {
		t.Errorf("round trip = %v, %v; want %v", got, err, due)
	}
	if formatInput(time.Time{}) != "" {
		t.Error("zero time should format as empty")
	}
}
