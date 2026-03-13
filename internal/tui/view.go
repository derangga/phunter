package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"phunter/internal/overlay"
)

func (m model) View() string {
	if m.quitting {
		return ""
	}

	helpBar := m.renderHelpBar()
	helpBarH := strings.Count(helpBar, "\n") + 1

	var autoRefreshBadge string
	if m.autoRefresh {
		autoRefreshBadge = m.styles.AutoRefreshOn.Render("● auto-refresh on")
	} else {
		autoRefreshBadge = m.styles.AutoRefreshOff.Render("○ auto-refresh off")
	}

	var sb strings.Builder
	sb.WriteString(m.styles.Title.Render("PortHunter: Hunt the active port and Kill it"))
	sb.WriteString("\n")
	sb.WriteString(autoRefreshBadge)
	sb.WriteString("\n")
	sb.WriteString(m.table.View())
	sb.WriteString("\n")
	sb.WriteString(m.styles.Status.Render(m.statusMsg))

	top := sb.String()

	// Fill remaining space so the help bar is pinned to the bottom.
	topLines := strings.Count(top, "\n") + 1
	gap := m.height - topLines - helpBarH
	if gap > 0 {
		top += strings.Repeat("\n", gap)
	}

	base := top + "\n" + helpBar

	if m.confirming && m.selected != nil && m.width > 0 {
		dialog := m.renderDialog()
		dw := lipgloss.Width(dialog)
		dh := strings.Count(dialog, "\n") + 1

		bgLines := strings.Split(base, "\n")
		bgH := len(bgLines)

		x := max((m.width-dw)/2, 0)
		y := max((bgH-helpBarH-dh)/2, 0)

		return overlay.Place(x, y, dialog, base)
	}

	return base
}

func (m model) renderDialog() string {
	question := m.styles.DialogBody.Render(
		fmt.Sprintf("Kill %s (%s) on :%s?", m.selected[1], m.selected[0], m.selected[5]),
	)
	hints := m.styles.Yes.Render("[y] Yes") + "   " + m.styles.No.Render("[N] No")
	return m.styles.DialogBox.Render(question + "\n\n" + hints)
}

func (m model) renderHelpBar() string {
	type binding struct{ key, desc string }
	bindings := []binding{
		{"↑/↓", "Navigate"},
		{"Enter/k", "Kill"},
		{"r", "Refresh"},
		{"a", "Auto-refresh"},
		{"q", "Quit"},
	}

	var parts []string
	for _, b := range bindings {
		parts = append(parts,
			m.styles.HelpKey.Render(b.key)+m.styles.HelpDesc.Render(b.desc),
		)
	}
	bar := strings.Join(parts, m.styles.HelpBar.Render(" "))
	// Pad the bar to full terminal width so the background fills the line.
	barW := lipgloss.Width(bar)
	if m.width > barW {
		bar += m.styles.HelpBar.Render(strings.Repeat(" ", m.width-barW))
	}
	return bar
}
