package tui

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"

	"phunter/internal/theme"
)

type Styles struct {
	Title     lipgloss.Style
	Status    lipgloss.Style
	HelpBar   lipgloss.Style
	HelpKey   lipgloss.Style
	HelpDesc  lipgloss.Style
	DialogBox lipgloss.Style
	DialogBody lipgloss.Style
	Yes       lipgloss.Style
	No        lipgloss.Style
	Table     table.Styles
}

func NewStyles(t theme.Theme) Styles {
	s := Styles{
		Title: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(t.Title)).
			MarginBottom(1),

		Status: lipgloss.NewStyle().
			Foreground(lipgloss.Color(t.StatusText)),

		HelpBar: lipgloss.NewStyle().
			Background(lipgloss.Color(t.HelpBarBg)).
			Foreground(lipgloss.Color(t.HelpBarFg)),

		HelpKey: lipgloss.NewStyle().
			Background(lipgloss.Color(t.HelpKeyBg)).
			Foreground(lipgloss.Color(t.HelpKeyFg)).
			Bold(true).
			Padding(0, 1),

		HelpDesc: lipgloss.NewStyle().
			Background(lipgloss.Color(t.HelpBarBg)).
			Foreground(lipgloss.Color(t.HelpBarFg)).
			Padding(0, 1),

		DialogBox: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(t.DialogBorder)).
			Padding(1, 2),

		DialogBody: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(t.DialogBody)),

		Yes: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(t.YesButton)),

		No: lipgloss.NewStyle().
			Foreground(lipgloss.Color(t.NoButton)),
	}

	ts := table.DefaultStyles()
	ts.Header = ts.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color(t.TableHeaderBorder)).
		BorderBottom(true).
		Bold(true)
	ts.Selected = ts.Selected.
		Foreground(lipgloss.Color(t.SelectedRowFg)).
		Background(lipgloss.Color(t.SelectedRowBg)).
		Bold(false)
	s.Table = ts

	return s
}
