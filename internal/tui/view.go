package tui

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"

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

	// 6. Help bar / confirm bar
	var footer string
	switch m.mode {
	case ModeConfirmKill:
		footer = m.renderConfirmBar()
	default:
		footer = m.renderHelpBar()
	}

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

	return topContent + "\n" + statusLine + "\n" + footer
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
	gap := m.width - leftW - rightW
	if gap < 2 {
		gap = 2
	}

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

	innerW := m.width - 6 // border + padding
	if innerW < 20 {
		innerW = 20
	}

	pillsW := lipgloss.Width(pills)
	inputW := innerW - pillsW - 3 // separator spacing
	if inputW < 10 {
		inputW = 10
	}
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
	end := m.offset + vh
	if end > len(m.viewProcs) {
		end = len(m.viewProcs)
	}

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

	// Apply highlighting
	if nameFilter != "" {
		nameStr = m.highlightCell(padOrTruncate(p.Name, nameW), nameFilter, nameW)
		userStr = m.highlightCell(padOrTruncate(p.User, colUser), nameFilter, colUser)
	}
	if portFilter != "" {
		portStr = m.highlightCell(padOrTruncate(p.Port, colPort), portFilter, colPort)
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
		left = m.styles.HeaderAccent.Render("▶") + " " +
			m.styles.StatusText.Render(p.Name) + "  " +
			m.styles.StatusDim.Render("pid") + " " +
			m.styles.StatusPid.Render(fmt.Sprintf("%d", p.PID)) + "  " +
			m.styles.StatusDim.Render("·  port") + " " +
			m.styles.StatusText.Render(p.Port) + "  " +
			m.styles.StatusDim.Render("·") + "  " +
			m.portClassStyle(class).Render(class.String())
	}

	// Right side: sort info + auto refresh countdown
	arrow := "↑"
	if !m.sortAsc {
		arrow = "↓"
	}
	right := m.styles.StatusDim.Render("sort") + " " +
		m.styles.StatusText.Bold(true).Render(m.sortKey.String()) + " " +
		m.styles.StatusText.Render(arrow)

	if m.autoRefresh {
		right += "  " + m.styles.StatusDim.Render("·") + "  " +
			m.styles.AutoRefreshOn.Render("●") + " " +
			m.styles.StatusDim.Render(fmt.Sprintf("auto in %ds", m.nextTickIn))
	}

	leftW := lipgloss.Width(left)
	rightW := lipgloss.Width(right)
	gap := m.width - leftW - rightW
	if gap < 2 {
		gap = 2
	}

	line := left + strings.Repeat(" ", gap) + right
	return m.styles.StatusBar.Width(m.width).Render(line)
}

// renderHelpBar renders the context-sensitive keybind footer.
func (m model) renderHelpBar() string {
	type binding struct{ key, desc string }

	var bindings []binding
	switch m.mode {
	case ModeFilter:
		bindings = []binding{
			{"Tab", "mode"},
			{"⏎", "apply"},
			{"Esc", "close"},
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

	var parts []string
	for _, b := range bindings {
		parts = append(parts,
			m.styles.HelpKey.Render(b.key)+m.styles.HelpDesc.Render(b.desc),
		)
	}
	bar := strings.Join(parts, m.styles.HelpBar.Render(" "))
	barW := lipgloss.Width(bar)
	if m.width > barW {
		bar += m.styles.HelpBar.Render(strings.Repeat(" ", m.width-barW))
	}
	return bar
}

// renderConfirmBar renders the kill confirmation footer.
func (m model) renderConfirmBar() string {
	// Find the target process name
	name := "?"
	for _, p := range m.viewProcs {
		if p.PID == m.killTarget {
			name = p.Name
			break
		}
	}

	question := m.styles.ConfirmBar.Render(
		fmt.Sprintf("Kill process %s (PID %d)?", name, m.killTarget))
	yes := m.styles.ConfirmYes.Render("[y] yes")
	no := m.styles.ConfirmNo.Render("[n] cancel")

	bar := question + "  " + yes + "  " + no
	barW := lipgloss.Width(bar)
	if m.width > barW {
		bar += m.styles.HelpBar.Render(strings.Repeat(" ", m.width-barW))
	}
	return bar
}

// Helper functions

func (m model) nameColWidth() int {
	fixed := colGlyph + colPID + colUser + colType + colAddr + colPort + 7 // 7 separators
	w := m.width - fixed
	if w < 12 {
		w = 12
	}
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
