package reminder

import (
	"fmt"
	"os/exec"
	"runtime"
)

// SystemNotifier shows a desktop notification using the OS's own tool.
// reminders.ENG.4
type SystemNotifier struct{}

func (SystemNotifier) Notify(title, body string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("notify-send", "--app-name=LoseNotion", title, body)
	case "darwin":
		script := fmt.Sprintf("display notification %q with title %q", body, title)
		cmd = exec.Command("osascript", "-e", script)
	default:
		return fmt.Errorf("notifications are not supported on %s", runtime.GOOS)
	}
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("%s: %w: %s", cmd.Path, err, out)
	}
	return nil
}
