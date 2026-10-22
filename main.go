package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/hridayK/LoseNotion/internal/reminder"
	"github.com/hridayK/LoseNotion/internal/store"
	"github.com/hridayK/LoseNotion/internal/tui"
)

const usage = `usage:
  losenotion                            open the todo list
  losenotion reminders [-interval 30s]  send reminder notifications (keep it running in the background)`

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "losenotion:", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	dir, err := store.DefaultDir()
	if err != nil {
		return err
	}
	st, err := store.Open(dir)
	if err != nil {
		return err
	}

	if len(args) == 0 {
		return runTUI(st)
	}
	switch args[0] {
	case "reminders":
		return runReminders(st, args[1:])
	case "help", "-h", "--help":
		fmt.Println(usage)
		return nil
	}
	return fmt.Errorf("unknown command %q\n%s", args[0], usage)
}

func runTUI(st *store.Store) error {
	// persistence.SAVE.3: stop before the UI can save over data we couldn't read
	list, err := st.Load()
	if err != nil {
		return err
	}
	_, err = tea.NewProgram(tui.New(list, st.Save, time.Now)).Run()
	return err
}

// reminders.ENG.2
func runReminders(st *store.Store, args []string) error {
	flags := flag.NewFlagSet("reminders", flag.ContinueOnError)
	interval := flags.Duration("interval", 30*time.Second, "how often to check for due reminders")
	if err := flags.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return nil
		}
		return err
	}
	if *interval <= 0 {
		return fmt.Errorf("-interval must be positive, got %s", *interval)
	}

	// Stop cleanly on Ctrl+C or when systemd stops the service.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	fmt.Fprintf(os.Stderr, "losenotion: checking reminders every %s (data in %s)\n", *interval, st.Dir())
	checker := reminder.Checker{Store: st, Notifier: reminder.SystemNotifier{}, Now: time.Now}
	return checker.Run(ctx, *interval, os.Stderr)
}
