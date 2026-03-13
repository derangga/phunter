package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"phunter/internal/theme"
	"phunter/internal/tui"
)

func main() {
	t := theme.Load()
	p := tea.NewProgram(tui.New(t), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
