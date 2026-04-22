package tui

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"

	"phunter/internal/overlay"
	"phunter/internal/ports"
	"phunter/internal/process"
)

func (m model) View() string {
	if m.quitting {
		return ""
	}

	var sections []string

	// 1. Header
	sections = append(sections, m.renderHeader())

	// 2. Filter bar (if active)
	if m.mode == ModeFilter || m.nameInput.Value() != "" || m.portInput.Value() != "" {
		sections = append(sections, m.renderFilterBar())
	}

	// 3. Table (header + rows via bubbles table)
	sections = append(sections, m.table.View())

	// 4. Status line
	statusLine := m.renderStatusLine()

	// 5. Footer help bar
	footer := m.renderHelpBar()

	// Calculate gap to push status + footer to bottom
	topContent := strings.Join(sections, "\n")
	topLines := strings.Count(topContent, "\n") + 1
	statusLines := 1
	footerLines := strings.Count(footer, "\n") + 1
	usedLines := topLines + statusLines + footerLines

	gap := m.height - usedLines
	if gap > 0 {
		topContent += strings.Repeat("\n", gap)
	}

	output := topContent + "\n" + statusLine + "\n" + footer

	// 6. Kill confirmation overlay
	if m.mode == ModeConfirmKill {
		dialog := m.renderKillDialog()
		dialogW := lipgloss.Width(dialog)
		dialogH := strings.Count(dialog, "\n") + 1
		x := (m.width - dialogW) / 2
		y := (m.height - dialogH) / 2
		output = overlay.Place(x, y, dialog, output)
	}

	// 7. Help overlay
	if m.showHelp {
		help := m.renderHelpOverlay()
		helpW := lipgloss.Width(help)
		helpH := strings.Count(help, "\n") + 1
		x := (m.width - helpW) / 2
		y := (m.height - helpH) / 2
		output = overlay.Place(x, y, help, output)
	}

	return output
}

// renderHeader renders: ◉ PortHunter  v0.3.0    listening N/M processes  │  ● auto  HH:MM:SS
func (m model) renderHeader() string {
	left := m.styles.HeaderAccent.Render(" ◉") + " " +
		m.styles.HeaderTitle.Render("PortHunter") + "  " +
		m.styles.HeaderDim.Render(m.version)

	total := len(m.allProcs)
	filtered := len(m.viewProcs)
	countStr := fmt.Sprintf("listening %s/%d processes",
		m.styles.HeaderAccent.Render(fmt.Sprintf("%d", filtered)),
		total)

	var autoStr string
	if m.autoRefresh {
		autoStr = m.styles.AutoRefreshOn.Render("● auto")
	} else {
		autoStr = m.styles.AutoRefreshOff.Render("○ auto")
	}

	right := countStr + "  " +
		m.styles.HeaderDim.Render("│") + "  " +
		autoStr + "  "

	leftW := lipgloss.Width(left)
	rightW := lipgloss.Width(right)
	gap := max(m.width-leftW-rightW, 2)

	return left + strings.Repeat(" ", gap) + right
}

// renderFilterBar renders the filter input bar with NAME|PORT mode pills.
func (m model) renderFilterBar() string {
	var namePill, portPill string
	if m.filterField == FilterName {
		namePill = m.styles.FilterModeActive.Render("NAME")
		portPill = m.styles.FilterModeIdle.Render("PORT")
	} else {
		namePill = m.styles.FilterModeIdle.Render("NAME")
		portPill = m.styles.FilterModeActive.Render("PORT")
	}

	pills := namePill + portPill

	var inputView string
	if m.filterField == FilterName {
		inputView = m.nameInput.View()
	} else {
		inputView = m.portInput.View()
	}

	innerW := max(m.width-6, 20)
	pillsW := lipgloss.Width(pills)
	inputW := max(innerW-pillsW-3, 10)
	m.nameInput.Width = inputW
	m.portInput.Width = inputW

	content := pills + "  " + m.styles.HeaderDim.Render("│") + " " + inputView

	bar := m.styles.FilterBorder.Width(m.width - 2).Render(content)
	return bar
}

// buildRows builds the table rows with plain text (no ANSI styling).
func (m model) buildRows() []table.Row {
	rows := make([]table.Row, len(m.viewProcs))
	for i, p := range m.viewProcs {
		rows[i] = m.buildRow(p)
	}
	return rows
}

func (m model) buildRow(p process.Process) table.Row {
	// Glyph cell: ● if multi-selected, else port class glyph (plain text only)
	var selectedCell string
	if m.selectedPIDs[p.PID] {
		selectedCell = "●"
	} else {
		selectedCell = ""
	}

	// Process cell: ● prefix if multi-selected
	nameCell := p.Name

	return table.Row{
		selectedCell,
		strconv.Itoa(p.PID),
		nameCell,
		p.User,
		p.Type,
		p.Address,
		p.Port,
	}
}

// renderStatusLine renders the status line between table and footer.
func (m model) renderStatusLine() string {
	var left string

	if m.toast != "" {
		left = m.styles.ToastStyle.Render(m.toast)
	} else if p := m.selectedProc(); p != nil {
		portNum, _ := strconv.Atoi(p.Port)
		class := ports.Classify(portNum)
		classStyle := m.portClassStyle(class)
		left = m.styles.HeaderAccent.Render("▸") + " " +
			m.styles.StatusText.Render(p.Name) + "  " +
			m.styles.StatusDim.Render("pid") + " " +
			m.styles.StatusPid.Render(fmt.Sprintf("%d", p.PID)) + "  " +
			m.styles.StatusDim.Render("·  port") + " " +
			classStyle.Render(p.Port) + "  " +
			m.styles.StatusDim.Render("·") + "  " +
			classStyle.Render(class.String())
	} else {
		left = m.styles.StatusDim.Italic(true).Render("no selection")
	}

	// Right side: filter indicators + sort info + auto refresh countdown
	var rightParts []string

	if v := m.nameInput.Value(); v != "" {
		rightParts = append(rightParts,
			m.styles.StatusDim.Render("name=")+m.styles.HeaderAccent.Render(v))
	}
	if v := m.portInput.Value(); v != "" {
		rightParts = append(rightParts,
			m.styles.StatusDim.Render("port=")+m.styles.HeaderAccent.Render(v))
	}

	arrow := "↑"
	if !m.sortAsc {
		arrow = "↓"
	}
	rightParts = append(rightParts,
		m.styles.StatusDim.Render("sort")+" "+
			m.styles.StatusText.Bold(true).Render(m.sortKey.String())+" "+
			m.styles.StatusText.Render(arrow))

	if m.autoRefresh {
		rightParts = append(rightParts,
			m.styles.AutoRefreshOn.Render("●")+" "+
				m.styles.StatusDim.Render(fmt.Sprintf("auto in %ds", m.nextTickIn)))
	} else {
		rightParts = append(rightParts,
			m.styles.StatusDim.Render("○ auto off"))
	}

	right := strings.Join(rightParts, "  "+m.styles.StatusDim.Render("·")+"  ")

	innerW := m.width - 2
	leftW := lipgloss.Width(left)
	rightW := lipgloss.Width(right)
	gap := max(innerW-leftW-rightW, 2)

	line := " " + left + strings.Repeat(" ", gap) + right + " "
	return m.styles.StatusBar.Width(m.width).Render(line)
}

// renderHelpBar renders the slim footer with only the most important keys.
func (m model) renderHelpBar() string {
	type binding struct{ key, desc string }

	var bindings []binding
	switch m.mode {
	case ModeFilter:
		bindings = []binding{
			{"⇥", "toggle"},
			{"Esc", "close"},
			{"⏎", "apply"},
		}
	default:
		bindings = []binding{
			{"↑↓", "nav"},
			{"Space", "select"},
			{"a", "auto"},
			{"⏎", "kill"},
			{"?", "help"},
		}
	}

	var chips []string
	for _, b := range bindings {
		chips = append(chips,
			m.styles.FooterKey.Render(b.key)+" "+m.styles.FooterLabel.Render(b.desc),
		)
	}

	// Width-aware wrapping
	var lines []string
	var currentLine string
	for i, chip := range chips {
		candidate := currentLine
		if candidate != "" {
			candidate += "   " + chip
		} else {
			candidate = chip
		}
		if currentLine != "" && lipgloss.Width(candidate) > m.width && m.width > 0 {
			lines = append(lines, currentLine)
			currentLine = chip
		} else {
			if i > 0 && currentLine != "" {
				currentLine += "   " + chip
			} else {
				currentLine = chip
			}
		}
	}
	if currentLine != "" {
		lines = append(lines, currentLine)
	}

	border := m.styles.TableHeaderBorder.Render(strings.Repeat("─", m.width))
	return border + "\n" + strings.Join(lines, "\n")
}

// renderKillDialog renders a centered overlay dialog for kill confirmation.
func (m model) renderKillDialog() string {
	var question string
	if len(m.selectedPIDs) > 0 {
		question = m.styles.DialogBody.Render(
			fmt.Sprintf("Kill %d selected process(es)?", len(m.selectedPIDs)))
	} else {
		name := "?"
		for _, p := range m.viewProcs {
			if p.PID == m.killTarget {
				name = p.Name
				break
			}
		}
		question = m.styles.DialogBody.Render(
			fmt.Sprintf("Kill %s (PID %d)?", name, m.killTarget))
	}

	yes := m.styles.ConfirmYes.Render("[y] yes")
	no := m.styles.ConfirmNo.Render("[n] cancel")

	content := question + "\n\n" + yes + "  " + no
	return m.styles.DialogBox.Render(content)
}

// renderHelpOverlay renders a floating overlay showing all keybindings.
func (m model) renderHelpOverlay() string {
	title := m.styles.HeaderTitle.Render("Keybindings")

	type binding struct{ key, desc string }
	bindings := []binding{
		{"↑ / ↓", "navigate"},
		{"Space", "select / deselect"},
		{"Enter / k", "kill selected"},
		{"/", "open filter"},
		{"s", "cycle sort"},
		{"r", "refresh"},
		{"a", "toggle auto-refresh"},
		{"?", "toggle this help"},
		{"q", "quit"},
	}

	var lines []string
	keyStyle := m.styles.FooterKey
	descStyle := m.styles.FooterLabel
	for _, b := range bindings {
		line := keyStyle.Render(b.key) + "  " + descStyle.Render(b.desc)
		lines = append(lines, line)
	}

	content := title + "\n\n" + strings.Join(lines, "\n")
	return m.styles.HelpOverlay.Render(content)
}

func (m model) portClassStyle(c ports.Class) lipgloss.Style {
	switch c {
	case ports.ClassPrivileged:
		return m.styles.PortPrivileged
	case ports.ClassDev:
		return m.styles.PortDev
	case ports.ClassRegistered:
		return m.styles.PortRegistered
	case ports.ClassEphemeral:
		return m.styles.PortEphemeral
	case ports.ClassAny:
		return m.styles.PortAny
	}
	return m.styles.PortAny
}
