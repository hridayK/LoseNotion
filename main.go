package main

import (
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	p := tea.NewProgram(
		newSimplePage("This app is under construction"),
	)
	if _, err := p.Run(); err != nil {
		panic(err)
	}
}
