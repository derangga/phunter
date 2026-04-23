package tui

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"

	"phunter/internal/theme"
)

type Styles struct {
	// Header
	HeaderTitle    lipgloss.Style
	HeaderDim      lipgloss.Style
	HeaderAccent   lipgloss.Style
	AutoRefreshOn  lipgloss.Style
	AutoRefreshOff lipgloss.Style

	// Table header
	TableHeaderDim    lipgloss.Style
	TableHeaderActive lipgloss.Style
	TableHeaderBorder lipgloss.Style

	// Port class glyph colors
	PortPrivileged lipgloss.Style
	PortDev        lipgloss.Style
	PortRegistered lipgloss.Style
	PortEphemeral  lipgloss.Style
	PortAny        lipgloss.Style

	// Filter bar
	FilterBorder     lipgloss.Style
	FilterModeActive lipgloss.Style
	FilterModeIdle   lipgloss.Style

	// Status line
	StatusBar  lipgloss.Style
	StatusPid  lipgloss.Style
	StatusText lipgloss.Style
	StatusDim  lipgloss.Style

	// Footer keybind bar
	FooterKey   lipgloss.Style
	FooterLabel lipgloss.Style
	Description lipgloss.Style

	// Kill confirmation
	ConfirmYes lipgloss.Style
	ConfirmNo  lipgloss.Style
	DialogBox  lipgloss.Style
	DialogBody lipgloss.Style

	// Toast
	ToastStyle lipgloss.Style

	// Bubbles table
	Table    table.Styles
	TableBox lipgloss.Style

	// Help overlay
	HelpOverlay lipgloss.Style
}

func NewStyles(t theme.Theme) Styles {

	s := Styles{
		// Header
		HeaderTitle: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(t.Accent)),
		HeaderDim: lipgloss.NewStyle().
			Foreground(lipgloss.Color(t.HeaderDimFg)),
		HeaderAccent: lipgloss.NewStyle().
			Foreground(lipgloss.Color(t.Title)).
			Bold(true),
		AutoRefreshOn: lipgloss.NewStyle().
			Foreground(lipgloss.Color(t.AutoRefreshOn)).
			Bold(true),
		AutoRefreshOff: lipgloss.NewStyle().
			Foreground(lipgloss.Color(t.AutoRefreshOff)),

		// Table header
		TableHeaderDim: lipgloss.NewStyle().
			Foreground(lipgloss.Color(t.HeaderDimFg)).
			Bold(true),
		TableHeaderActive: lipgloss.NewStyle().
			Foreground(lipgloss.Color(t.SortActiveFg)).
			Bold(true),
		TableHeaderBorder: lipgloss.NewStyle().
			Foreground(lipgloss.Color(t.TableHeaderBorder)),

		// Port class
		PortPrivileged: lipgloss.NewStyle().
			Background(lipgloss.Color(t.StatusBarBg)).
			Foreground(lipgloss.Color(t.PortPrivileged)),
		PortDev: lipgloss.NewStyle().
			Background(lipgloss.Color(t.StatusBarBg)).
			Foreground(lipgloss.Color(t.PortDev)),
		PortRegistered: lipgloss.NewStyle().
			Background(lipgloss.Color(t.StatusBarBg)).
			Foreground(lipgloss.Color(t.PortRegistered)),
		PortEphemeral: lipgloss.NewStyle().
			Background(lipgloss.Color(t.StatusBarBg)).
			Foreground(lipgloss.Color(t.PortEphemeral)),
		PortAny: lipgloss.NewStyle().
			Background(lipgloss.Color(t.StatusBarBg)).
			Foreground(lipgloss.Color(t.PortAny)),

		// Filter bar
		FilterBorder: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(t.FilterBorder)),
		FilterModeActive: lipgloss.NewStyle().
			Background(lipgloss.Color(t.FilterModeActiveBg)).
			Foreground(lipgloss.Color(t.FilterModeActiveFg)).
			Bold(true).
			Padding(0, 1),
		FilterModeIdle: lipgloss.NewStyle().
			Background(lipgloss.Color(t.FilterModeIdleBg)).
			Foreground(lipgloss.Color(t.FilterModeIdleFg)).
			Padding(0, 1),

		// Status line
		StatusBar: lipgloss.NewStyle().
			Background(lipgloss.Color(t.StatusBarBg)),
		StatusPid: lipgloss.NewStyle().
			Background(lipgloss.Color(t.StatusBarBg)).
			Foreground(lipgloss.Color(t.StatusPidFg)).
			Bold(true),
		StatusText: lipgloss.NewStyle().
			Background(lipgloss.Color(t.StatusBarBg)).
			Foreground(lipgloss.Color(t.StatusText)),
		StatusDim: lipgloss.NewStyle().
			Background(lipgloss.Color(t.StatusBarBg)).
			Foreground(lipgloss.Color(t.HeaderDimFg)),

		// Footer keybind bar
		FooterKey: lipgloss.NewStyle().
			Background(lipgloss.Color(t.FooterKeyBg)).
			Foreground(lipgloss.Color(t.FooterKeyFg)).
			Bold(true).
			Padding(0, 1),
		FooterLabel: lipgloss.NewStyle().
			Foreground(lipgloss.Color(t.FooterLabelFg)),
		Description: lipgloss.NewStyle().
			Foreground(lipgloss.Color(t.Description)),

		// Kill confirmation
		ConfirmYes: lipgloss.NewStyle().
			Background(lipgloss.Color(t.ConfirmChipBg)).
			Foreground(lipgloss.Color(t.YesButton)).
			Bold(true).
			Padding(0, 1),
		ConfirmNo: lipgloss.NewStyle().
			Background(lipgloss.Color(t.ConfirmChipBg)).
			Foreground(lipgloss.Color(t.NoButton)).
			Padding(0, 1),
		DialogBox: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(t.DialogBorder)).
			Padding(1, 2),
		DialogBody: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(t.DialogBody)),

		// Toast
		ToastStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color(t.YesButton)),

		// Help overlay
		HelpOverlay: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(t.DialogBorder)).
			Padding(1, 2),
	}

	// Bubbles table styles
	ts := table.DefaultStyles()
	ts.Header = ts.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color(t.TableHeaderBorder)).
		BorderBottom(true).
		Bold(true).
		Foreground(lipgloss.Color(t.HeaderDimFg))
	ts.Selected = ts.Selected.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color(t.TableHeaderBorder)).
		BorderLeft(true).
		BorderRight(true).
		Foreground(lipgloss.Color(t.SelectedRowFg)).
		Background(lipgloss.Color(t.SelectedRowBg)).
		Bold(false)
	ts.Cell = ts.Cell.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color(t.TableHeaderBorder)).
		BorderLeft(false).
		BorderRight(false).
		BorderTop(false).
		BorderBottom(false)
	ts.Selected = ts.Selected.
		BorderLeft(false).
		BorderRight(false).
		BorderTop(false).
		BorderBottom(false).
		Foreground(lipgloss.Color(t.SelectedRowFg)).
		Background(lipgloss.Color(t.SelectedRowBg)).
		Bold(false)
	s.Table = ts
	s.TableBox = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color(t.TableHeaderBorder)).
		Padding(0, 0)

	return s
}
