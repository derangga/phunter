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

[colors]
title              = "#8aadf4"   # Blue
status_text        = "#a5adcb"   # Subtext0
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
filter_mode_active_bg = "#8aadf4" # accent
filter_mode_active_fg = "#24273a" # Base
filter_mode_idle_bg = "#363a4f"   # Surface0
filter_mode_idle_fg = "#a5adcb"   # Subtext0

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
confirm_chip_bg     = "#24273a"   # Base — chip background for kill confirm
confirm_chip_fg     = "#a6da95"   # Red — chip foreground for kill confirm

# Header / footer
header_dim_fg       = "#6e738d"   # Overlay0
footer_key_bg       = "#363a4f"   # Surface0 — key chip background
footer_key_fg       = "#b7bdf8"   # Lavender — key chip foreground
footer_label_fg     = "#8087a2"   # Overlay1 — key label foreground
`

type Theme struct {
	Title             string `toml:"title"`
	StatusText        string `toml:"status_text"`
	DialogBorder      string `toml:"dialog_border"`
	DialogBody        string `toml:"dialog_body"`
	YesButton         string `toml:"yes_button"`
	NoButton          string `toml:"no_button"`
	TableHeaderBorder string `toml:"table_header_border"`
	SelectedRowFg     string `toml:"selected_row_fg"`
	SelectedRowBg     string `toml:"selected_row_bg"`
	AutoRefreshOn     string `toml:"auto_refresh_on"`
	AutoRefreshOff    string `toml:"auto_refresh_off"`

	Accent             string `toml:"accent"`
	FilterBorder       string `toml:"filter_border"`
	FilterModeActiveBg string `toml:"filter_mode_active_bg"`
	FilterModeActiveFg string `toml:"filter_mode_active_fg"`
	FilterModeIdleBg   string `toml:"filter_mode_idle_bg"`
	FilterModeIdleFg   string `toml:"filter_mode_idle_fg"`
	SortActiveFg       string `toml:"sort_active_fg"`
	PortPrivileged     string `toml:"port_privileged"`
	PortDev            string `toml:"port_dev"`
	PortRegistered     string `toml:"port_registered"`
	PortEphemeral      string `toml:"port_ephemeral"`
	PortAny            string `toml:"port_any"`
	StatusBarBg        string `toml:"status_bar_bg"`
	StatusPidFg        string `toml:"status_pid_fg"`
	ConfirmBarBg       string `toml:"confirm_bar_bg"`
	ConfirmChipBg      string `toml:"confirm_chip_bg"`
	ConfirmChipFg      string `toml:"confirm_chip_fg"`
	HeaderDimFg        string `toml:"header_dim_fg"`
	FooterKeyBg        string `toml:"footer_key_bg"`
	FooterKeyFg        string `toml:"footer_key_fg"`
	FooterLabelFg      string `toml:"footer_label_fg"`
}

type configFile struct {
	Colors Theme `toml:"colors"`
}

func DefaultTheme() Theme {
	return Theme{
		Title:              "#8aadf4",
		StatusText:         "#a5adcb",
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
		FilterModeActiveBg: "#8aadf4",
		FilterModeActiveFg: "#24273a",
		FilterModeIdleBg:   "#363a4f",
		FilterModeIdleFg:   "#a5adcb",
		SortActiveFg:       "#8aadf4",
		PortPrivileged:     "#f5a97b",
		PortDev:            "#a6da95",
		PortRegistered:     "#7dc4e4",
		PortEphemeral:      "#6e738d",
		PortAny:            "#6e738d",
		StatusBarBg:        "#1e2030",
		StatusPidFg:        "#f5a97b",
		ConfirmBarBg:       "#ed8796",
		ConfirmChipBg:      "#24273a",
		ConfirmChipFg:      "#a6da95",
		HeaderDimFg:        "#6e738d",
		FooterKeyBg:        "#363a4f",
		FooterKeyFg:        "#b7bdf8",
		FooterLabelFg:      "#8087a2",
	}
}

// Load reads the theme from the XDG config file. If missing or invalid, returns defaults.
func Load() Theme {
	path, err := xdg.SearchConfigFile(configPath)
	if err != nil {
		return DefaultTheme()
	}

	var cfg configFile
	if _, err := toml.DecodeFile(path, &cfg); err != nil {
		fmt.Fprintf(os.Stderr, "phunter: warning: failed to parse %s: %v (using defaults)\n", path, err)
		return DefaultTheme()
	}
	return cfg.Colors
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
