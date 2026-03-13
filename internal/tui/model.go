package tui

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"

	"phunter/internal/process"
	"phunter/internal/theme"
)

const autoRefreshInterval = 5 * time.Second

type refreshMsg struct {
	processes []process.Process
	err       error
}

type tickMsg time.Time

type model struct {
	table      table.Model
	styles     Styles
	processes  []process.Process
	statusMsg  string
	confirming bool
	selected   table.Row
	quitting    bool
	autoRefresh bool
	width       int
	height      int
}

// New creates and returns the initial TUI model.
func New(th theme.Theme) model {
	styles := NewStyles(th)

	t := table.New(
		table.WithColumns(tableColumns(20)),
		table.WithFocused(true),
		table.WithHeight(20),
	)
	t.SetStyles(styles.Table)

	return model{table: t, styles: styles}
}

func refreshCmd() tea.Cmd {
	return func() tea.Msg {
		procs, err := process.GetListeningPorts()
		return refreshMsg{processes: procs, err: err}
	}
}

func tickCmd() tea.Cmd {
	return tea.Tick(autoRefreshInterval, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m model) Init() tea.Cmd {
	return refreshCmd()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		nameW := max(m.width-colPID-colUser-colType-colAddr-colPort, 12)
		m.table.SetColumns(tableColumns(nameW))
		tableH := max(m.height-7, 1)
		m.table.SetHeight(tableH)
		return m, nil

	case tickMsg:
		if m.autoRefresh {
			return m, tea.Batch(refreshCmd(), tickCmd())
		}
		return m, nil

	case refreshMsg:
		if msg.err != nil {
			m.statusMsg = "Error: " + msg.err.Error()
			return m, nil
		}
		m.processes = msg.processes
		rows := make([]table.Row, len(m.processes))
		for i, p := range m.processes {
			rows[i] = table.Row{
				strconv.Itoa(p.PID),
				p.Name,
				p.User,
				p.Type,
				p.Address,
				p.Port,
			}
		}
		m.table.SetRows(rows)
		if m.statusMsg == "" || strings.HasPrefix(m.statusMsg, "Refreshing") || strings.HasPrefix(m.statusMsg, "Auto-refresh") {
			m.statusMsg = fmt.Sprintf("%d process(es) listening", len(m.processes))
		}
		return m, nil

	case tea.KeyMsg:
		if m.confirming {
			switch msg.String() {
			case "y":
				return m.killSelected()
			case "n", "esc":
				m.confirming = false
				m.statusMsg = "Cancelled"
			}
			return m, nil
		}

		switch msg.String() {
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		case "r":
			m.statusMsg = "Refreshing..."
			return m, refreshCmd()
		case "a":
			m.autoRefresh = !m.autoRefresh
			if m.autoRefresh {
				m.statusMsg = "Auto-refresh ON (5s)"
				return m, tickCmd()
			}
			m.statusMsg = "Auto-refresh OFF"
			return m, nil
		case "enter", "k":
			if len(m.processes) == 0 {
				return m, nil
			}
			m.selected = m.table.SelectedRow()
			if m.selected == nil {
				return m, nil
			}
			m.confirming = true
			return m, nil
		}
	}

	// Delegate navigation keys to the table when browsing.
	if !m.confirming {
		var cmd tea.Cmd
		m.table, cmd = m.table.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m model) killSelected() (tea.Model, tea.Cmd) {
	m.confirming = false
	if m.selected == nil {
		return m, nil
	}
	pid, err := strconv.Atoi(m.selected[0])
	if err != nil {
		m.statusMsg = "Invalid PID"
		return m, nil
	}
	if err := process.Kill(pid); err != nil {
		m.statusMsg = err.Error()
		return m, nil
	}
	m.statusMsg = fmt.Sprintf("Killed %s (PID %d) on :%s", m.selected[1], pid, m.selected[5])
	return m, refreshCmd()
}
