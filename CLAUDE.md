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

<!-- gitnexus:start -->
# GitNexus — Code Intelligence

This project is indexed by GitNexus as **portkiller-tui** (372 symbols, 734 relationships, 23 execution flows). Use the GitNexus MCP tools to understand code, assess impact, and navigate safely.

> If any GitNexus tool warns the index is stale, run `npx gitnexus analyze` in terminal first.

## Always Do

- **MUST run impact analysis before editing any symbol.** Before modifying a function, class, or method, run `gitnexus_impact({target: "symbolName", direction: "upstream"})` and report the blast radius (direct callers, affected processes, risk level) to the user.
- **MUST run `gitnexus_detect_changes()` before committing** to verify your changes only affect expected symbols and execution flows.
- **MUST warn the user** if impact analysis returns HIGH or CRITICAL risk before proceeding with edits.
- When exploring unfamiliar code, use `gitnexus_query({query: "concept"})` to find execution flows instead of grepping. It returns process-grouped results ranked by relevance.
- When you need full context on a specific symbol — callers, callees, which execution flows it participates in — use `gitnexus_context({name: "symbolName"})`.

## When Debugging

1. `gitnexus_query({query: "<error or symptom>"})` — find execution flows related to the issue
2. `gitnexus_context({name: "<suspect function>"})` — see all callers, callees, and process participation
3. `READ gitnexus://repo/portkiller-tui/process/{processName}` — trace the full execution flow step by step
4. For regressions: `gitnexus_detect_changes({scope: "compare", base_ref: "main"})` — see what your branch changed

## When Refactoring

- **Renaming**: MUST use `gitnexus_rename({symbol_name: "old", new_name: "new", dry_run: true})` first. Review the preview — graph edits are safe, text_search edits need manual review. Then run with `dry_run: false`.
- **Extracting/Splitting**: MUST run `gitnexus_context({name: "target"})` to see all incoming/outgoing refs, then `gitnexus_impact({target: "target", direction: "upstream"})` to find all external callers before moving code.
- After any refactor: run `gitnexus_detect_changes({scope: "all"})` to verify only expected files changed.

## Never Do

- NEVER edit a function, class, or method without first running `gitnexus_impact` on it.
- NEVER ignore HIGH or CRITICAL risk warnings from impact analysis.
- NEVER rename symbols with find-and-replace — use `gitnexus_rename` which understands the call graph.
- NEVER commit changes without running `gitnexus_detect_changes()` to check affected scope.

## Tools Quick Reference

| Tool | When to use | Command |
|------|-------------|---------|
| `query` | Find code by concept | `gitnexus_query({query: "auth validation"})` |
| `context` | 360-degree view of one symbol | `gitnexus_context({name: "validateUser"})` |
| `impact` | Blast radius before editing | `gitnexus_impact({target: "X", direction: "upstream"})` |
| `detect_changes` | Pre-commit scope check | `gitnexus_detect_changes({scope: "staged"})` |
| `rename` | Safe multi-file rename | `gitnexus_rename({symbol_name: "old", new_name: "new", dry_run: true})` |
| `cypher` | Custom graph queries | `gitnexus_cypher({query: "MATCH ..."})` |

## Impact Risk Levels

| Depth | Meaning | Action |
|-------|---------|--------|
| d=1 | WILL BREAK — direct callers/importers | MUST update these |
| d=2 | LIKELY AFFECTED — indirect deps | Should test |
| d=3 | MAY NEED TESTING — transitive | Test if critical path |

## Resources

| Resource | Use for |
|----------|---------|
| `gitnexus://repo/portkiller-tui/context` | Codebase overview, check index freshness |
| `gitnexus://repo/portkiller-tui/clusters` | All functional areas |
| `gitnexus://repo/portkiller-tui/processes` | All execution flows |
| `gitnexus://repo/portkiller-tui/process/{name}` | Step-by-step execution trace |

## Self-Check Before Finishing

Before completing any code modification task, verify:
1. `gitnexus_impact` was run for all modified symbols
2. No HIGH/CRITICAL risk warnings were ignored
3. `gitnexus_detect_changes()` confirms changes match expected scope
4. All d=1 (WILL BREAK) dependents were updated

## Keeping the Index Fresh

After committing code changes, the GitNexus index becomes stale. Re-run analyze to update it:

```bash
npx gitnexus analyze
```

If the index previously included embeddings, preserve them by adding `--embeddings`:

```bash
npx gitnexus analyze --embeddings
```

To check whether embeddings exist, inspect `.gitnexus/meta.json` — the `stats.embeddings` field shows the count (0 means no embeddings). **Running analyze without `--embeddings` will delete any previously generated embeddings.**

> Claude Code users: A PostToolUse hook handles this automatically after `git commit` and `git merge`.

## CLI

| Task | Read this skill file |
|------|---------------------|
| Understand architecture / "How does X work?" | `.claude/skills/gitnexus/gitnexus-exploring/SKILL.md` |
| Blast radius / "What breaks if I change X?" | `.claude/skills/gitnexus/gitnexus-impact-analysis/SKILL.md` |
| Trace bugs / "Why is X failing?" | `.claude/skills/gitnexus/gitnexus-debugging/SKILL.md` |
| Rename / extract / split / refactor | `.claude/skills/gitnexus/gitnexus-refactoring/SKILL.md` |
| Tools, resources, schema reference | `.claude/skills/gitnexus/gitnexus-guide/SKILL.md` |
| Index, status, clean, wiki CLI commands | `.claude/skills/gitnexus/gitnexus-cli/SKILL.md` |

<!-- gitnexus:end -->
