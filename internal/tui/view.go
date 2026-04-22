package tui

import (
	"fmt"
	"strconv"
	"strings"
	"time"

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

	// 3. Table header
	sections = append(sections, m.renderTableHeader())

	// 4. Table rows
	tableContent := m.renderTableRows()

	// 5. Status line
	statusLine := m.renderStatusLine()

	// 6. Footer (always help bar)
	footer := m.renderHelpBar()

	// Calculate gap to push status + footer to bottom
	topContent := strings.Join(sections, "\n") + "\n" + tableContent
	topLines := strings.Count(topContent, "\n") + 1
	statusLines := 1
	footerLines := strings.Count(footer, "\n") + 1
	usedLines := topLines + statusLines + footerLines

	gap := m.height - usedLines
	if gap > 0 {
		topContent += strings.Repeat("\n", gap)
	}

	output := topContent + "\n" + statusLine + "\n" + footer

	// 7. Kill confirmation overlay
	if m.mode == ModeConfirmKill {
		dialog := m.renderKillDialog()
		dialogW := lipgloss.Width(dialog)
		dialogH := strings.Count(dialog, "\n") + 1
		x := (m.width - dialogW) / 2
		y := (m.height - dialogH) / 2
		output = overlay.Place(x, y, dialog, output)
	}

	return output
}

// renderHeader renders: ◉ PortHunter  v0.3.0    listening N/M processes  │  ● auto  HH:MM:SS
func (m model) renderHeader() string {
	left := m.styles.HeaderAccent.Render("◉") + " " +
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

	clock := m.styles.HeaderDim.Render(time.Now().Format("15:04:05"))

	right := countStr + "  " +
		m.styles.HeaderDim.Render("│") + "  " +
		autoStr + "  " + clock

	// Pad middle to fill width
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

	// border + padding
	innerW := max(m.width-6, 20)

	pillsW := lipgloss.Width(pills)

	// separator spacing
	inputW := max(innerW-pillsW-3, 10)
	m.nameInput.Width = inputW
	m.portInput.Width = inputW

	content := pills + "  " + m.styles.HeaderDim.Render("│") + " " + inputView

	bar := m.styles.FilterBorder.Width(m.width - 2).Render(content)
	return bar
}

// renderTableHeader renders ALL-CAPS column headers with sort indicator.
func (m model) renderTableHeader() string {
	nameW := m.nameColWidth()

	type col struct {
		title string
		width int
		key   SortKey
	}

	cols := []col{
		{"PID", colPID, SortPID},
		{"PROCESS", nameW, SortProcess},
		{"USER", colUser, SortUser},
		{"TYPE", colType, SortType},
		{"ADDRESS", colAddr, -1}, // not sortable, no SortKey
		{"PORT", colPort, SortPort},
	}

	var parts []string
	// Glyph column spacer
	parts = append(parts, strings.Repeat(" ", colGlyph))

	for _, c := range cols {
		title := c.title
		style := m.styles.TableHeaderDim

		if c.key >= 0 && c.key == m.sortKey {
			arrow := " ↑"
			if !m.sortAsc {
				arrow = " ↓"
			}
			title += arrow
			style = m.styles.TableHeaderActive
		}

		rendered := style.Render(padOrTruncate(title, c.width))
		parts = append(parts, rendered)
	}

	header := strings.Join(parts, " ")

	// Border line under header
	border := m.styles.TableHeaderBorder.Render(strings.Repeat("─", m.width))

	return header + "\n" + border
}

// renderTableRows renders the visible rows.
func (m model) renderTableRows() string {
	if len(m.viewProcs) == 0 {
		if m.nameInput.Value() != "" || m.portInput.Value() != "" {
			return m.styles.EmptyState.Render("no matches. press esc to clear filter.")
		}
		return m.styles.EmptyState.Render("no listening processes found.")
	}

	vh := m.viewHeight()
	end := min(m.offset+vh, len(m.viewProcs))

	nameW := m.nameColWidth()
	nameFilter := m.nameInput.Value()
	portFilter := m.portInput.Value()

	var rows []string
	for i := m.offset; i < end; i++ {
		p := m.viewProcs[i]
		isSelected := i == m.cursor

		row := m.renderRow(p, isSelected, nameW, nameFilter, portFilter)
		rows = append(rows, row)
	}

	return strings.Join(rows, "\n")
}

func (m model) renderRow(p process.Process, selected bool, nameW int, nameFilter, portFilter string) string {
	portNum, _ := strconv.Atoi(p.Port)
	class := ports.Classify(portNum)

	// Glyph column
	glyphStyle := m.portClassStyle(class)
	var glyph string
	if selected && m.selStyle == "bar" {
		glyph = m.styles.RowSelectedBar.Render("▎") + glyphStyle.Render(class.Glyph())
	} else {
		glyph = " " + glyphStyle.Render(class.Glyph())
	}

	// Build cell values
	pidStr := padOrTruncate(strconv.Itoa(p.PID), colPID)
	nameStr := padOrTruncate(p.Name, nameW)
	userStr := padOrTruncate(p.User, colUser)
	typeStr := padOrTruncate(p.Type, colType)
	addrStr := padOrTruncate(p.Address, colAddr)
	portStr := padOrTruncate(p.Port, colPort)

	// Apply TYPE column coloring (IPv4 = sapphire, IPv6 = mauve)
	if strings.Contains(p.Type, "6") {
		typeStr = m.styles.TypeIPv6.Render(typeStr)
	} else {
		typeStr = m.styles.TypeIPv4.Render(typeStr)
	}

	// Apply highlighting
	if nameFilter != "" {
		nameStr = m.highlightCell(padOrTruncate(p.Name, nameW), nameFilter, nameW)
		userStr = m.highlightCell(padOrTruncate(p.User, colUser), nameFilter, colUser)
	}
	if portFilter != "" {
		portStr = m.highlightCell(padOrTruncate(p.Port, colPort), portFilter, colPort)
	}

	// Bold selected process name
	if selected {
		nameStr = lipgloss.NewStyle().Bold(true).Render(nameStr)
	}

	// Port number gets class color
	if portFilter == "" {
		portStr = glyphStyle.Render(padOrTruncate(p.Port, colPort))
	}

	cells := glyph + " " + pidStr + " " + nameStr + " " + userStr + " " + typeStr + " " + addrStr + " " + portStr

	// Apply row-level style
	if selected {
		if m.selStyle == "block" {
			return m.styles.RowSelectedBlock.Width(m.width).Render(cells)
		}
		return m.styles.RowSelected.Width(m.width).Render(cells)
	}
	return m.styles.RowNormal.Width(m.width).Render(cells)
}

func (m model) highlightCell(text, filter string, width int) string {
	before, match, after := highlightMatch(text, filter)
	if match == "" {
		return text
	}
	result := before + m.styles.MatchHighlight.Render(match) + after
	return result
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

	// Account for 1-char padding on each side
	innerW := m.width - 2
	leftW := lipgloss.Width(left)
	rightW := lipgloss.Width(right)
	gap := max(innerW-leftW-rightW, 2)

	line := " " + left + strings.Repeat(" ", gap) + right + " "
	return m.styles.StatusBar.Width(m.width).Render(line)
}

// renderHelpBar renders the context-sensitive keybind footer.
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
			{"/", "filter"},
			{"s", "sort"},
			{"⏎", "kill"},
			{"r", "refresh"},
			{"a", "auto"},
			{"g/G", "top/end"},
			{"q", "quit"},
		}
	}

	// Build chips: key chip (surface0 bg + lavender fg) + label (overlay1 fg)
	var chips []string
	for _, b := range bindings {
		chips = append(chips,
			m.styles.FooterKey.Render(b.key)+" "+m.styles.FooterLabel.Render(b.desc),
		)
	}

	// Width-aware wrapping: greedy line-break
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

	// Top border + bar (no solid background — transparent)
	border := m.styles.TableHeaderBorder.Render(strings.Repeat("─", m.width))
	return border + "\n" + strings.Join(lines, "\n")
}

// renderConfirmBar renders the kill confirmation footer.
func (m model) renderConfirmBar() string {
	question := m.styles.ConfirmBar.Render(
		fmt.Sprintf("▲  Kill process %d?", m.killTarget))
	yes := m.styles.ConfirmYes.Render("[y] yes")
	no := m.styles.ConfirmNo.Render("[n] cancel")

	right := yes + "  " + no
	leftW := lipgloss.Width(question)
	rightW := lipgloss.Width(right)
	gap := max(m.width-leftW-rightW, 2)

	bar := question + m.styles.ConfirmBar.Render(strings.Repeat(" ", gap)) + right
	barW := lipgloss.Width(bar)
	if m.width > barW {
		bar += m.styles.ConfirmBar.Render(strings.Repeat(" ", m.width-barW))
	}

	border := m.styles.TableHeaderBorder.Render(strings.Repeat("─", m.width))
	return border + "\n" + bar
}

// renderKillDialog renders a centered overlay dialog for kill confirmation.
func (m model) renderKillDialog() string {
	// Find the target process name
	name := "?"
	for _, p := range m.viewProcs {
		if p.PID == m.killTarget {
			name = p.Name
			break
		}
	}

	question := m.styles.DialogBody.Render(
		fmt.Sprintf("Kill %s (PID %d)?", name, m.killTarget))

	yes := m.styles.ConfirmYes.Render("[y] yes")
	no := m.styles.ConfirmNo.Render("[n] cancel")

	buttons := yes + "  " + no
	content := question + "\n\n" + buttons

	return m.styles.DialogBox.Render(content)
}

// Helper functions

func (m model) nameColWidth() int {
	fixed := colGlyph + colPID + colUser + colType + colAddr + colPort + 7 // 7 separators
	w := max(m.width-fixed, 12)
	return w
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

func padOrTruncate(s string, width int) string {
	if width <= 0 {
		return ""
	}
	runes := []rune(s)
	if len(runes) > width {
		if width > 1 {
			return string(runes[:width-1]) + "…"
		}
		return string(runes[:width])
	}
	return s + strings.Repeat(" ", width-len(runes))
}
