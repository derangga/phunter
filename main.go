package main

import (
	"fmt"
	"os"

	"phunter/internal/theme"
	"phunter/internal/tui"

	tea "github.com/charmbracelet/bubbletea"
)

// version is set at build time via -ldflags "-X main.version=..."
var version = "dev"

func main() {
	if len(os.Args) == 2 && (os.Args[1] == "-v" || os.Args[1] == "--version") {
		fmt.Printf("phunter %s\n", version)
		return
	}

	theme.EnsureConfig()
	t := theme.Load()
	p := tea.NewProgram(tui.New(t, version), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
