package theme

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/adrg/xdg"
)

const configPath = "phunter/config.toml"

const defaultConfig = `# phunter Theme Configuration
# Edit colors using hex values (e.g. "#f5bde6")
# Delete this file to regenerate with defaults.

# Selection style: "bar" (left-rail accent) or "block" (full-width inverted row)
# selection_style = "bar"

[colors]
title              = "#8aadf4"   # Blue
status_text        = "#a5adcb"   # Subtext0
help_bar_bg        = "#363a4f"   # Surface0
help_bar_fg        = "#b8c0e0"   # Subtext1
help_key_bg        = "#89b4fa"   # Blue
help_key_fg        = "#181926"   # Crust
dialog_border      = "#89b4fa"   # Blue
dialog_body        = "#cad3f5"   # Text
yes_button         = "#a6da95"   # Green
no_button          = "#8087a2"   # Overlay1
table_header_border = "#494d64"  # Surface1
selected_row_fg    = "#181926"   # Crust
selected_row_bg    = "#89b4fa"   # Blue
auto_refresh_on    = "#a6da95"   # Green
auto_refresh_off   = "#6e738d"   # Overlay0

# --- New optional keys (v0.3.0+) ---

# Accent color for filter bar focus, sort arrows, match highlights
accent              = "#8aadf4"   # Blue

# Filter bar
filter_border       = "#363a4f"   # Surface0
filter_label_fg     = "#a5adcb"   # Subtext0
filter_placeholder  = "#6e738d"   # Overlay0
filter_mode_active_bg = "#8aadf4" # accent
filter_mode_active_fg = "#24273a" # Base
filter_mode_idle_bg = "#363a4f"   # Surface0
filter_mode_idle_fg = "#a5adcb"   # Subtext0
match_highlight_bg  = "#8aadf4"   # accent
match_highlight_fg  = "#24273a"   # Base

# Sort arrow / active column header
sort_active_fg      = "#8aadf4"   # accent

# Port class glyphs
port_privileged     = "#f5a97b"   # Peach
port_dev            = "#a6da95"   # Green
port_registered     = "#7dc4e4"   # Sapphire
port_ephemeral      = "#6e738d"   # Overlay0
port_any            = "#6e738d"   # Overlay0

# Status line
status_bar_bg       = "#1e2030"   # Mantle
status_pid_fg       = "#f5a97b"   # Peach
confirm_bar_bg      = "#ed8796"   # Red — background for kill confirm footer
confirm_bar_fg      = "#24273a"   # Base — foreground for kill confirm
confirm_chip_bg     = "#24273a"   # Base — inverted chip background
confirm_chip_fg     = "#ed8796"   # Red — inverted chip foreground

# Row selection
row_selected_bar    = "#8aadf4"   # accent — left rail color

# Row TYPE column
row_type_ipv4_fg    = "#7dc4e4"   # Sapphire — IPv4 text color
row_type_ipv6_fg    = "#c6a0f6"   # Mauve — IPv6 text color

# Header
header_fg           = "#cad3f5"   # Text
header_dim_fg       = "#6e738d"   # Overlay0

# Footer keybind bar
footer_key_bg       = "#363a4f"   # Surface0 — key chip background
footer_key_fg       = "#b7bdf8"   # Lavender — key chip foreground
footer_label_fg     = "#8087a2"   # Overlay1 — key label foreground
`

type Theme struct {
	// Existing keys
	Title             string `toml:"title"`
	StatusText        string `toml:"status_text"`
	HelpBarBg         string `toml:"help_bar_bg"`
	HelpBarFg         string `toml:"help_bar_fg"`
	HelpKeyBg         string `toml:"help_key_bg"`
	HelpKeyFg         string `toml:"help_key_fg"`
	DialogBorder      string `toml:"dialog_border"`
	DialogBody        string `toml:"dialog_body"`
	YesButton         string `toml:"yes_button"`
	NoButton          string `toml:"no_button"`
	TableHeaderBorder string `toml:"table_header_border"`
	SelectedRowFg     string `toml:"selected_row_fg"`
	SelectedRowBg     string `toml:"selected_row_bg"`
	AutoRefreshOn     string `toml:"auto_refresh_on"`
	AutoRefreshOff    string `toml:"auto_refresh_off"`

	// New optional keys (v0.3.0+)
	Accent             string `toml:"accent"`
	FilterBorder       string `toml:"filter_border"`
	FilterLabelFg      string `toml:"filter_label_fg"`
	FilterPlaceholder  string `toml:"filter_placeholder"`
	FilterModeActiveBg string `toml:"filter_mode_active_bg"`
	FilterModeActiveFg string `toml:"filter_mode_active_fg"`
	FilterModeIdleBg   string `toml:"filter_mode_idle_bg"`
	FilterModeIdleFg   string `toml:"filter_mode_idle_fg"`
	MatchHighlightBg   string `toml:"match_highlight_bg"`
	MatchHighlightFg   string `toml:"match_highlight_fg"`
	SortActiveFg       string `toml:"sort_active_fg"`
	PortPrivileged     string `toml:"port_privileged"`
	PortDev            string `toml:"port_dev"`
	PortRegistered     string `toml:"port_registered"`
	PortEphemeral      string `toml:"port_ephemeral"`
	PortAny            string `toml:"port_any"`
	StatusBarBg        string `toml:"status_bar_bg"`
	StatusPidFg        string `toml:"status_pid_fg"`
	ConfirmBarBg       string `toml:"confirm_bar_bg"`
	ConfirmBarFg       string `toml:"confirm_bar_fg"`
	ConfirmChipBg      string `toml:"confirm_chip_bg"`
	ConfirmChipFg      string `toml:"confirm_chip_fg"`
	RowSelectedBar     string `toml:"row_selected_bar"`
	HeaderFg           string `toml:"header_fg"`
	HeaderDimFg        string `toml:"header_dim_fg"`
	FooterKeyBg        string `toml:"footer_key_bg"`
	FooterKeyFg        string `toml:"footer_key_fg"`
	FooterLabelFg      string `toml:"footer_label_fg"`
	RowTypeIPv4Fg      string `toml:"row_type_ipv4_fg"`
	RowTypeIPv6Fg      string `toml:"row_type_ipv6_fg"`
}

type configFile struct {
	SelectionStyle string `toml:"selection_style"`
	Colors         Theme  `toml:"colors"`
}

// Config holds the full parsed configuration including non-color settings.
type Config struct {
	Theme          Theme
	SelectionStyle string // "bar" or "block"
}

func DefaultTheme() Theme {
	return Theme{
		Title:             "#8aadf4",
		StatusText:         "#a5adcb",
		HelpBarBg:          "#363a4f",
		HelpBarFg:          "#b8c0e0",
		HelpKeyBg:          "#c6a0f6",
		HelpKeyFg:          "#181926",
		DialogBorder:       "#8aadf4",
		DialogBody:         "#cad3f5",
		YesButton:          "#a6da95",
		NoButton:           "#8087a2",
		TableHeaderBorder:  "#494d64",
		SelectedRowFg:      "#181926",
		SelectedRowBg:      "#c6a0f6",
		AutoRefreshOn:      "#a6da95",
		AutoRefreshOff:     "#6e738d",
		Accent:             "#8aadf4",
		FilterBorder:       "#363a4f",
		FilterLabelFg:      "#a5adcb",
		FilterPlaceholder:  "#6e738d",
		FilterModeActiveBg: "#8aadf4",
		FilterModeActiveFg: "#24273a",
		FilterModeIdleBg:   "#363a4f",
		FilterModeIdleFg:   "#a5adcb",
		MatchHighlightBg:   "#8aadf4",
		MatchHighlightFg:   "#24273a",
		SortActiveFg:       "#8aadf4",
		PortPrivileged:     "#f5a97b",
		PortDev:            "#a6da95",
		PortRegistered:     "#7dc4e4",
		PortEphemeral:      "#6e738d",
		PortAny:            "#6e738d",
		StatusBarBg:        "#1e2030",
		StatusPidFg:        "#f5a97b",
		ConfirmBarBg:       "#ed8796",
		ConfirmBarFg:       "#24273a",
		ConfirmChipBg:      "#24273a",
		ConfirmChipFg:      "#ed8796",
		RowSelectedBar:     "#8aadf4",
		HeaderFg:           "#cad3f5",
		HeaderDimFg:        "#6e738d",
		FooterKeyBg:        "#363a4f",
		FooterKeyFg:        "#b7bdf8",
		FooterLabelFg:      "#8087a2",
		RowTypeIPv4Fg:      "#7dc4e4",
		RowTypeIPv6Fg:      "#c6a0f6",
	}
}

// resolve fills zero-value new fields with sensible fallbacks from existing keys.
func (t *Theme) resolve() {
	def := DefaultTheme()
	fallback := func(val *string, primary, secondary string) {
		if *val == "" {
			if primary != "" {
				*val = primary
			} else {
				*val = secondary
			}
		}
	}

	fallback(&t.AutoRefreshOn, t.YesButton, def.AutoRefreshOn)
	fallback(&t.AutoRefreshOff, "", def.AutoRefreshOff)
	fallback(&t.Accent, t.HelpKeyBg, def.Accent)
	fallback(&t.FilterBorder, t.TableHeaderBorder, def.FilterBorder)
	fallback(&t.FilterLabelFg, t.StatusText, def.FilterLabelFg)
	fallback(&t.FilterPlaceholder, "", def.FilterPlaceholder)
	fallback(&t.FilterModeActiveBg, t.Accent, def.FilterModeActiveBg)
	fallback(&t.FilterModeActiveFg, "", def.FilterModeActiveFg)
	fallback(&t.FilterModeIdleBg, t.HelpBarBg, def.FilterModeIdleBg)
	fallback(&t.FilterModeIdleFg, t.StatusText, def.FilterModeIdleFg)
	fallback(&t.MatchHighlightBg, t.Accent, def.MatchHighlightBg)
	fallback(&t.MatchHighlightFg, "", def.MatchHighlightFg)
	fallback(&t.SortActiveFg, t.Accent, def.SortActiveFg)
	fallback(&t.PortPrivileged, "", def.PortPrivileged)
	fallback(&t.PortDev, "", def.PortDev)
	fallback(&t.PortRegistered, "", def.PortRegistered)
	fallback(&t.PortEphemeral, "", def.PortEphemeral)
	fallback(&t.PortAny, "", def.PortAny)
	fallback(&t.StatusBarBg, "", def.StatusBarBg)
	fallback(&t.StatusPidFg, "", def.StatusPidFg)
	fallback(&t.ConfirmBarBg, t.DialogBorder, def.ConfirmBarBg)
	fallback(&t.ConfirmBarFg, "", def.ConfirmBarFg)
	fallback(&t.ConfirmChipBg, "", def.ConfirmChipBg)
	fallback(&t.ConfirmChipFg, t.ConfirmBarBg, def.ConfirmChipFg)
	fallback(&t.RowSelectedBar, t.Accent, def.RowSelectedBar)
	fallback(&t.HeaderFg, t.DialogBody, def.HeaderFg)
	fallback(&t.HeaderDimFg, "", def.HeaderDimFg)
	fallback(&t.FooterKeyBg, t.HelpBarBg, def.FooterKeyBg)
	fallback(&t.FooterKeyFg, t.HelpBarFg, def.FooterKeyFg)
	fallback(&t.FooterLabelFg, t.HelpBarFg, def.FooterLabelFg)
	fallback(&t.RowTypeIPv4Fg, t.PortRegistered, def.RowTypeIPv4Fg)
	fallback(&t.RowTypeIPv6Fg, t.DialogBorder, def.RowTypeIPv6Fg)
}

// Load reads the theme from the XDG config file. If missing or invalid, returns defaults.
func Load() Config {
	path, err := xdg.SearchConfigFile(configPath)
	if err != nil {
		return Config{Theme: DefaultTheme(), SelectionStyle: "bar"}
	}

	var cfg configFile
	if _, err := toml.DecodeFile(path, &cfg); err != nil {
		fmt.Fprintf(os.Stderr, "phunter: warning: failed to parse %s: %v (using defaults)\n", path, err)
		return Config{Theme: DefaultTheme(), SelectionStyle: "bar"}
	}

	cfg.Colors.resolve()

	style := cfg.SelectionStyle
	if style == "" {
		style = "bar"
	}

	return Config{Theme: cfg.Colors, SelectionStyle: style}
}

// EnsureConfig writes the default config file if it doesn't already exist.
func EnsureConfig() {
	if _, err := xdg.SearchConfigFile(configPath); err == nil {
		return // already exists
	}
	path, err := xdg.ConfigFile(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "phunter: warning: could not create config path: %v\n", err)
		return
	}
	if err := os.WriteFile(path, []byte(defaultConfig), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "phunter: warning: could not write config: %v\n", err)
	}
}
