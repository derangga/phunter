# Code Style & Conventions

## Go Conventions
- Standard Go formatting (`go fmt`)
- All packages live under `internal/` — nothing is exported outside the module
- No tests currently exist in the project

## Architecture Patterns
- **Elm architecture**: Model → Update → View (via Bubble Tea)
- Styles are never hardcoded in view/model code; they flow from `theme.Theme` → `tui.Styles`
- `tui.New()` accepts a `theme.Theme`; `main.go` calls `theme.Load()` to provide it
- Table column widths are constants in `columns.go`; process name width is computed from terminal width
- Use `tea.Cmd` for async operations (refresh, quit); never block in `Update()`

## Naming
- Unexported types for internal model (`model`, `refreshMsg`, `tickMsg`)
- Exported types for cross-package use (`Process`, `Theme`, `Styles`)
- Constants use camelCase (`colPID`, `autoRefreshInterval`)

## Adding a New Theme Color
1. Add field to `Theme` struct and `DefaultTheme()` in `internal/theme/theme.go`
2. Add TOML key + default to `defaultConfig` const in same file
3. Add `lipgloss.Style` field to `Styles` in `internal/tui/styles.go`
4. Build the style from theme color in `NewStyles()`
5. Reference via `m.styles.YourField` in view code
