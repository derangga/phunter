package tui

import (
	"fmt"
	"strconv"
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

type clearStatusMsg struct{}

type killMsg struct {
	name string
	pid  int
	port string
	err  error
}

type batchKillMsg struct {
	killed int
	errors []string
}

type model struct {
	table        table.Model
	styles       Styles
	processes    []process.Process
	statusMsg    string
	statusLocked bool
	confirming   bool
	selected     *process.Process
	selectedPIDs map[int]bool
	quitting     bool
	autoRefresh  bool
	width        int
	height       int
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

	return model{table: t, styles: styles, selectedPIDs: make(map[int]bool)}
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

func clearStatusCmd() tea.Cmd {
	return tea.Tick(3*time.Second, func(time.Time) tea.Msg {
		return clearStatusMsg{}
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
		nameW := max(m.width-colPID-colUser-colType-colAddr-colPort-12, 12)
		m.table.SetColumns(tableColumns(nameW))
		tableH := max(m.height-7, 1)
		m.table.SetHeight(tableH)
		return m, nil

	case tickMsg:
		if m.autoRefresh {
			return m, tea.Batch(refreshCmd(), tickCmd())
		}
		return m, nil

	case batchKillMsg:
		m.selectedPIDs = make(map[int]bool)
		if len(msg.errors) > 0 {
			m.statusMsg = fmt.Sprintf("Killed %d, %d failed: %s", msg.killed, len(msg.errors), msg.errors[0])
		} else {
			m.statusMsg = fmt.Sprintf("Killed %d process(es)", msg.killed)
		}
		m.statusLocked = true
		return m, tea.Batch(refreshCmd(), clearStatusCmd())

	case killMsg:
		if msg.err != nil {
			m.statusMsg = msg.err.Error()
			m.statusLocked = true
			return m, clearStatusCmd()
		}
		m.statusMsg = fmt.Sprintf("Killed %s (PID %d) on :%s", msg.name, msg.pid, msg.port)
		m.statusLocked = true
		return m, tea.Batch(refreshCmd(), clearStatusCmd())

	case clearStatusMsg:
		m.statusLocked = false
		m.statusMsg = fmt.Sprintf("%d process(es) listening", len(m.processes))
		return m, nil

	case refreshMsg:
		if msg.err != nil {
			m.statusMsg = "Error: " + msg.err.Error()
			return m, nil
		}
		m.processes = msg.processes
		m.table.SetRows(m.buildRows())
		if !m.statusLocked {
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
				m.selectedPIDs = make(map[int]bool)
				m.table.SetRows(m.buildRows())
				m.statusMsg = "Cancelled"
				m.statusLocked = true
				return m, clearStatusCmd()
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
		case " ":
			if len(m.processes) == 0 {
				return m, nil
			}
			cursor := m.table.Cursor()
			if cursor < 0 || cursor >= len(m.processes) {
				return m, nil
			}
			pid := m.processes[cursor].PID
			if m.selectedPIDs[pid] {
				delete(m.selectedPIDs, pid)
			} else {
				m.selectedPIDs[pid] = true
			}
			m.table.SetRows(m.buildRows())
			if len(m.selectedPIDs) > 0 {
				m.statusMsg = fmt.Sprintf("%d selected", len(m.selectedPIDs))
			} else {
				m.statusMsg = fmt.Sprintf("%d process(es) listening", len(m.processes))
			}
			return m, nil
		case "enter", "k":
			if len(m.processes) == 0 {
				return m, nil
			}
			if len(m.selectedPIDs) > 0 {
				// Batch mode
				m.selected = nil
				m.confirming = true
				return m, nil
			}
			cursor := m.table.Cursor()
			if cursor < 0 || cursor >= len(m.processes) {
				return m, nil
			}
			m.selected = &m.processes[cursor]
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

func (m model) buildRows() []table.Row {
	rows := make([]table.Row, len(m.processes))
	for i, p := range m.processes {
		name := p.Name
		if m.selectedPIDs[p.PID] {
			name = "● " + name
		}
		rows[i] = table.Row{
			strconv.Itoa(p.PID),
			name,
			p.User,
			p.Type,
			p.Address,
			p.Port,
		}
	}
	return rows
}

func killCmd(p process.Process) tea.Cmd {
	return func() tea.Msg {
		err := process.Kill(p.PID)
		return killMsg{name: p.Name, pid: p.PID, port: p.Port, err: err}
	}
}

func batchKillCmd(procs []process.Process) tea.Cmd {
	return func() tea.Msg {
		var killed int
		var errs []string
		for _, p := range procs {
			if err := process.Kill(p.PID); err != nil {
				errs = append(errs, err.Error())
			} else {
				killed++
			}
		}
		return batchKillMsg{killed: killed, errors: errs}
	}
}

func (m model) killSelected() (tea.Model, tea.Cmd) {
	m.confirming = false

	// Batch kill mode
	if len(m.selectedPIDs) > 0 {
		var targets []process.Process
		for _, p := range m.processes {
			if m.selectedPIDs[p.PID] {
				targets = append(targets, p)
			}
		}
		m.statusMsg = fmt.Sprintf("Killing %d process(es)...", len(targets))
		return m, batchKillCmd(targets)
	}

	// Single kill mode
	if m.selected == nil {
		return m, nil
	}
	m.statusMsg = fmt.Sprintf("Killing %s (PID %d)...", m.selected.Name, m.selected.PID)
	return m, killCmd(*m.selected)
}
