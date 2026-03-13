# phunter

> Port Hunter — a terminal UI for finding and killing processes listening on TCP ports.

![demo placeholder](https://via.placeholder.com/800x400?text=demo+GIF+goes+here)

## Features

- Live table of all TCP `LISTEN` processes with PID, name, user, address, and port
- Kill any process with a confirmation dialog — no accidental kills
- Auto-refresh after every kill
- Fullscreen alt-screen TUI (no terminal pollution)
- Themeable via a simple TOML config file (Catppuccin Macchiato defaults)
- Zero runtime dependencies — single static binary

## Installation

### Build from source

```sh
git clone https://github.com/sociolla/portkiller-tui
cd portkiller-tui
go build -o phunter .
./phunter
```

### go install

```sh
go install phunter@latest
```

### Homebrew

```sh
# coming soon
brew install phunter
```

### Nix

```sh
nix run github:sociolla/portkiller-tui
```

## Usage

```sh
phunter
```

No flags or arguments. The table populates immediately with all listening TCP ports.

## Keybindings

### Browsing

| Key | Action |
|-----|--------|
| `↑` / `↓` | Navigate rows |
| `Enter` or `k` | Open kill confirmation dialog |
| `r` | Refresh process list |
| `q` / `Ctrl+C` | Quit |

### Confirmation dialog

| Key | Action |
|-----|--------|
| `y` | Confirm kill |
| `n` / `Esc` | Cancel |

## Table columns

| Column | Description |
|--------|-------------|
| PID | Process ID |
| Process | Executable name |
| User | Owner of the process |
| Type | Protocol (TCP4, TCP6, …) |
| Address | Bind address |
| Port | Listening port number |

## Configuration

The config file is created automatically on first run:

```
~/.config/phunter/config.toml
```

(follows XDG Base Directory; on macOS this resolves to `~/Library/Application Support/phunter/config.toml` when `$XDG_CONFIG_HOME` is unset, but the `~/.config` path is preferred if set)

Delete the file to regenerate it with defaults.

### Full config reference

```toml
# phunter Theme Configuration
# Edit colors using hex values (e.g. "#f5bde6")
# Delete this file to regenerate with defaults.

[colors]
title               = "#f5bde6"   # Pink         — app title
status_text         = "#a5adcb"   # Subtext0     — status / count text
help_bar_bg         = "#363a4f"   # Surface0     — help bar background
help_bar_fg         = "#b8c0e0"   # Subtext1     — help bar text
help_key_bg         = "#c6a0f6"   # Mauve        — key chip background
help_key_fg         = "#181926"   # Crust        — key chip text
dialog_border       = "#c6a0f6"   # Mauve        — kill dialog border
dialog_body         = "#cad3f5"   # Text         — kill dialog body text
yes_button          = "#a6da95"   # Green        — [y] confirm button
no_button           = "#8087a2"   # Overlay1     — [n] cancel button
table_header_border = "#494d64"   # Surface1     — line under table header
selected_row_fg     = "#181926"   # Crust        — highlighted row text
selected_row_bg     = "#c6a0f6"   # Mauve        — highlighted row background
```

All 13 color keys are required when overriding; any parse error falls back to the defaults above with a warning printed to stderr.

## Requirements

- **Go 1.21+** (module requires go 1.25)
- **`lsof`** available in `$PATH`
- **macOS** (amd64 / arm64) — primary target
- **Linux** (amd64 / arm64) — fully supported
- **Windows** (amd64) — partial support (`lsof` availability may vary)

## License

MIT
