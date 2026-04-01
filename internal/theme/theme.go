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
title              = "#f5bde6"   # Pink
status_text        = "#a5adcb"   # Subtext0
help_bar_bg        = "#363a4f"   # Surface0
help_bar_fg        = "#b8c0e0"   # Subtext1
help_key_bg        = "#c6a0f6"   # Mauve
help_key_fg        = "#181926"   # Crust
dialog_border      = "#c6a0f6"   # Mauve
dialog_body        = "#cad3f5"   # Text
yes_button         = "#a6da95"   # Green
no_button          = "#8087a2"   # Overlay1
table_header_border = "#494d64"  # Surface1
selected_row_fg    = "#181926"   # Crust
selected_row_bg    = "#c6a0f6"   # Mauve
auto_refresh_on    = "#a6da95"   # Green
auto_refresh_off   = "#6e738d"   # Overlay0
`

type Theme struct {
	Title            string `toml:"title"`
	StatusText       string `toml:"status_text"`
	HelpBarBg        string `toml:"help_bar_bg"`
	HelpBarFg        string `toml:"help_bar_fg"`
	HelpKeyBg        string `toml:"help_key_bg"`
	HelpKeyFg        string `toml:"help_key_fg"`
	DialogBorder     string `toml:"dialog_border"`
	DialogBody       string `toml:"dialog_body"`
	YesButton        string `toml:"yes_button"`
	NoButton         string `toml:"no_button"`
	TableHeaderBorder string `toml:"table_header_border"`
	SelectedRowFg    string `toml:"selected_row_fg"`
	SelectedRowBg    string `toml:"selected_row_bg"`
	AutoRefreshOn    string `toml:"auto_refresh_on"`
	AutoRefreshOff   string `toml:"auto_refresh_off"`
}

type configFile struct {
	Colors Theme `toml:"colors"`
}

func DefaultTheme() Theme {
	return Theme{
		Title:            "#f5bde6",
		StatusText:       "#a5adcb",
		HelpBarBg:        "#363a4f",
		HelpBarFg:        "#b8c0e0",
		HelpKeyBg:        "#c6a0f6",
		HelpKeyFg:        "#181926",
		DialogBorder:     "#c6a0f6",
		DialogBody:       "#cad3f5",
		YesButton:        "#a6da95",
		NoButton:         "#8087a2",
		TableHeaderBorder: "#494d64",
		SelectedRowFg:    "#181926",
		SelectedRowBg:    "#c6a0f6",
		AutoRefreshOn:    "#a6da95",
		AutoRefreshOff:   "#6e738d",
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
