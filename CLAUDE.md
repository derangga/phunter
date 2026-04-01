# phunter (Port Hunter)

Terminal UI for hunting and killing processes listening on TCP ports. Built with Go + Bubble Tea.

## Project Structure

```
main.go                        # Entry point: loads theme, runs Bubble Tea program
internal/
  tui/
    model.go                   # Model struct, Init(), Update(), key handling
    view.go                    # View(), renderDialog(), renderHelpBar()
    styles.go                  # Styles struct built from theme colors
    columns.go                 # Table column width constants
  process/process.go           # GetListeningPorts() via lsof, Kill() via SIGKILL
  theme/theme.go               # XDG config loading, Catppuccin Macchiato defaults
  overlay/overlay.go           # Dialog overlay placement with ANSI preservation
```

## Stack

- **Go** with **Bubble Tea** (bubbletea), **Bubbles** (table component), **Lipgloss** (styling)
- **BurntSushi/toml** for config parsing, **adrg/xdg** for cross-platform config paths
- Binary name: `phunter`

## Architecture

- Single-binary TUI app using the Elm architecture (Model → Update → View)
- Two UI states: **browsing** (table navigation) and **confirming** (kill dialog overlay)
- Process list from `lsof -i -n -P -sTCP:LISTEN`, deduplicated by `pid:port`
- Kill sends `SIGKILL` via `os.FindProcess()` then auto-refreshes
- Theme colors loaded from `phunter/config.toml` in XDG config dir on startup
- Config auto-generated on first run; parse errors fall back to defaults with stderr warning

## Build & Run

```sh
go build -o phunter .
./phunter
```

## Conventions

- All packages live under `internal/` — nothing is exported outside the module
- Styles are never hardcoded in view/model code; they flow from `theme.Theme` → `tui.Styles`
- `tui.New()` accepts a `theme.Theme`; `main.go` calls `theme.Load()` to provide it
- Table column widths are constants in `columns.go`; process name width is computed from terminal width
- Use `tea.Cmd` for async operations (refresh, quit); never block in `Update()`

## Adding a New Theme Color

1. Add the field to `Theme` struct and `DefaultTheme()` in `internal/theme/theme.go`
2. Add the TOML key + default to the `defaultConfig` const in the same file
3. Add a `lipgloss.Style` field to `Styles` in `internal/tui/styles.go`
4. Build the style from the theme color in `NewStyles()`
5. Reference via `m.styles.YourField` in view code

## Serena MCP

When the Serena MCP server is available, you **MUST** use Serena tools instead of raw file tools. This is a **BLOCKING REQUIREMENT** — do NOT use `Read`, `Edit`, `Grep`, or `Glob` on source files when Serena can accomplish the task.

### Required tool mapping

| Instead of | MUST use |
|---|---|
| `Read` a source file | `get_symbols_overview` → `find_symbol` (read only what you need) |
| `Grep` / `Glob` for code | `search_for_pattern` or `find_symbol` with substring matching |
| `Edit` source code | `replace_symbol_body`, `insert_before_symbol`, `insert_after_symbol` |
| Navigating references | `find_referencing_symbols` |

### Exceptions (raw tools allowed)

- Non-source files (go.mod, go.sum, .toml, .yml, .md, .gitignore)
- When a Serena tool explicitly fails or cannot handle the operation
- Creating entirely new files (`Write` tool)

## RTK (Rust Token Killer)

All shell commands must be prefixed with `rtk` to route through the token-optimized CLI proxy:

```sh
rtk go build -o phunter .
rtk go test ./...
rtk go fmt ./...
rtk go vet ./...
rtk go mod tidy
rtk git status
rtk git diff
rtk git log
```

**Exception**: `rtk` meta commands (`rtk gain`, `rtk discover`, `rtk proxy`) are used directly without a subcommand prefix.

## Bubbletea Skill

When working on TUI-related tasks (components, views, styling, key handling, layout), **always invoke the `bubbletea` skill first** to get up-to-date patterns and best practices before writing or modifying TUI code.
