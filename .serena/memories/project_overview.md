# phunter (Port Hunter) — Project Overview

## Purpose
Terminal UI application for hunting and killing processes listening on TCP ports. Users can browse listening ports in a table, select a process, and kill it with a confirmation dialog.

## Tech Stack
- **Language**: Go (1.25.7)
- **TUI Framework**: Bubble Tea (bubbletea v1.3.10) — Elm architecture (Model → Update → View)
- **Components**: Bubbles v1.0.0 (table component)
- **Styling**: Lipgloss v1.1.0
- **Config**: BurntSushi/toml for parsing, adrg/xdg for cross-platform config paths
- **Binary name**: `phunter`
- **Module name**: `phunter`

## Architecture
- Single-binary TUI app with two UI states: **browsing** (table navigation) and **confirming** (kill dialog overlay)
- Process list obtained from `lsof -i -n -P -sTCP:LISTEN`, deduplicated by `pid:port`
- Kill sends `SIGKILL` via `os.FindProcess()` then auto-refreshes
- Theme colors loaded from `phunter/config.toml` in XDG config dir; auto-generated on first run
- Uses `tea.WithAltScreen()` for fullscreen TUI

## Codebase Structure
```
main.go                          # Entry point: loads theme, runs Bubble Tea program
internal/
  tui/
    model.go                     # Model struct, Init(), Update(), key handling, New(), refreshCmd, tickCmd
    view.go                      # View(), renderDialog(), renderHelpBar()
    styles.go                    # Styles struct built from theme colors via NewStyles()
    columns.go                   # Table column width constants (colPID, colUser, colType, colAddr, colPort), tableColumns()
  process/process.go             # Process struct, GetListeningPorts() via lsof, Kill() via SIGKILL
  theme/theme.go                 # Theme struct, DefaultTheme(), Load(), generate(); XDG config, Catppuccin Macchiato defaults
  overlay/overlay.go             # Place(), takeColumns(), skipColumns() — dialog overlay with ANSI preservation
```

## Release & Distribution
- GoReleaser v2 (`.goreleaser.yaml`)
- Platforms: linux/darwin (amd64, arm64), windows (amd64)
- Homebrew tap: `derangga/homebrew-formulae`
- Nix flake: `flake.nix`
- Packages: .deb, .rpm via nfpms
- CI: GitHub Actions (`.github/workflows/release.yml`) on tag push
