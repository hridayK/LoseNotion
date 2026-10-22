package reminder

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestSystemNotifier(t *testing.T) {
	t.Run("reminders.ENG.4 uses notify-send on Linux", func(t *testing.T) {
		if runtime.GOOS != "linux" {
			t.Skip("Linux only")
		}
		// Put a fake notify-send first on PATH that records its arguments.
		dir := t.TempDir()
		argsFile := filepath.Join(dir, "args")
		script := "#!/bin/sh\nprintf '%s\\n' \"$@\" > " + argsFile + "\n"
		if err := os.WriteFile(filepath.Join(dir, "notify-send"), []byte(script), 0o755); err != nil {
			t.Fatal(err)
		}
		t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))

		if err := (SystemNotifier{}).Notify("LoseNotion reminder", "buy milk"); err != nil {
			t.Fatal(err)
		}
		got, _ := os.ReadFile(argsFile)
		want := "--app-name=LoseNotion\nLoseNotion reminder\nbuy milk\n"
		if string(got) != want {
			t.Errorf("notify-send args = %q, want %q", got, want)
		}
	})
}
