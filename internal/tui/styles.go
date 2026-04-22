package tui

import (
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

	// Rows
	RowNormal        lipgloss.Style
	RowSelected      lipgloss.Style
	RowSelectedBar   lipgloss.Style
	RowSelectedBlock lipgloss.Style
	TypeIPv4         lipgloss.Style
	TypeIPv6         lipgloss.Style

	// Port class glyph colors
	PortPrivileged lipgloss.Style
	PortDev        lipgloss.Style
	PortRegistered lipgloss.Style
	PortEphemeral  lipgloss.Style
	PortAny        lipgloss.Style

	// Filter bar
	FilterBorder      lipgloss.Style
	FilterLabel       lipgloss.Style
	FilterPlaceholder lipgloss.Style
	FilterModeActive  lipgloss.Style
	FilterModeIdle    lipgloss.Style
	MatchHighlight    lipgloss.Style

	// Status line
	StatusBar   lipgloss.Style
	StatusPid   lipgloss.Style
	StatusText  lipgloss.Style
	StatusDim   lipgloss.Style

	// Help bar (legacy)
	HelpBar  lipgloss.Style
	HelpKey  lipgloss.Style
	HelpDesc lipgloss.Style

	// Footer keybind bar
	FooterKey   lipgloss.Style
	FooterLabel lipgloss.Style

	// Kill confirmation
	ConfirmBar  lipgloss.Style
	ConfirmYes  lipgloss.Style
	ConfirmNo   lipgloss.Style
	DialogBox   lipgloss.Style
	DialogBody  lipgloss.Style

	// Toast
	ToastStyle lipgloss.Style

	// Empty state
	EmptyState lipgloss.Style
}

func NewStyles(t theme.Theme) Styles {
	accent := lipgloss.Color(t.Accent)
	surface0 := lipgloss.Color(t.HelpBarBg)

	s := Styles{
		// Header
		HeaderTitle: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(t.Title)),
		HeaderDim: lipgloss.NewStyle().
			Foreground(lipgloss.Color(t.HeaderDimFg)),
		HeaderAccent: lipgloss.NewStyle().
			Foreground(accent).
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

		// Rows
		RowNormal: lipgloss.NewStyle(),
		RowSelected: lipgloss.NewStyle().
			Background(surface0),
		RowSelectedBar: lipgloss.NewStyle().
			Foreground(accent),
		RowSelectedBlock: lipgloss.NewStyle().
			Foreground(lipgloss.Color(t.SelectedRowFg)).
			Background(lipgloss.Color(t.SelectedRowBg)),
		TypeIPv4: lipgloss.NewStyle().
			Foreground(lipgloss.Color(t.RowTypeIPv4Fg)),
		TypeIPv6: lipgloss.NewStyle().
			Foreground(lipgloss.Color(t.RowTypeIPv6Fg)),

		// Port class
		PortPrivileged: lipgloss.NewStyle().
			Foreground(lipgloss.Color(t.PortPrivileged)),
		PortDev: lipgloss.NewStyle().
			Foreground(lipgloss.Color(t.PortDev)),
		PortRegistered: lipgloss.NewStyle().
			Foreground(lipgloss.Color(t.PortRegistered)),
		PortEphemeral: lipgloss.NewStyle().
			Foreground(lipgloss.Color(t.PortEphemeral)),
		PortAny: lipgloss.NewStyle().
			Foreground(lipgloss.Color(t.PortAny)),

		// Filter bar
		FilterBorder: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(t.FilterBorder)),
		FilterLabel: lipgloss.NewStyle().
			Foreground(lipgloss.Color(t.FilterLabelFg)),
		FilterPlaceholder: lipgloss.NewStyle().
			Foreground(lipgloss.Color(t.FilterPlaceholder)),
		FilterModeActive: lipgloss.NewStyle().
			Background(lipgloss.Color(t.FilterModeActiveBg)).
			Foreground(lipgloss.Color(t.FilterModeActiveFg)).
			Bold(true).
			Padding(0, 1),
		FilterModeIdle: lipgloss.NewStyle().
			Background(lipgloss.Color(t.FilterModeIdleBg)).
			Foreground(lipgloss.Color(t.FilterModeIdleFg)).
			Padding(0, 1),
		MatchHighlight: lipgloss.NewStyle().
			Background(lipgloss.Color(t.MatchHighlightBg)).
			Foreground(lipgloss.Color(t.MatchHighlightFg)),

		// Status line
		StatusBar: lipgloss.NewStyle().
			Background(lipgloss.Color(t.StatusBarBg)),
		StatusPid: lipgloss.NewStyle().
			Foreground(lipgloss.Color(t.StatusPidFg)).
			Bold(true),
		StatusText: lipgloss.NewStyle().
			Foreground(lipgloss.Color(t.StatusText)),
		StatusDim: lipgloss.NewStyle().
			Foreground(lipgloss.Color(t.HeaderDimFg)),

		// Help bar (legacy, kept for backward compat)
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

		// Footer keybind bar (new design)
		FooterKey: lipgloss.NewStyle().
			Background(lipgloss.Color(t.FooterKeyBg)).
			Foreground(lipgloss.Color(t.FooterKeyFg)).
			Bold(true).
			Padding(0, 1),
		FooterLabel: lipgloss.NewStyle().
			Foreground(lipgloss.Color(t.FooterLabelFg)),

		// Kill confirmation
		ConfirmBar: lipgloss.NewStyle().
			Background(lipgloss.Color(t.ConfirmBarBg)).
			Foreground(lipgloss.Color(t.ConfirmBarFg)),
		ConfirmYes: lipgloss.NewStyle().
			Background(lipgloss.Color(t.ConfirmChipBg)).
			Foreground(lipgloss.Color(t.ConfirmChipFg)).
			Bold(true).
			Padding(0, 1),
		ConfirmNo: lipgloss.NewStyle().
			Background(lipgloss.Color(t.ConfirmChipBg)).
			Foreground(lipgloss.Color(t.ConfirmChipFg)).
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

		// Empty state
		EmptyState: lipgloss.NewStyle().
			Foreground(lipgloss.Color(t.HeaderDimFg)).
			Italic(true).
			Padding(1, 2),
	}

	return s
}
