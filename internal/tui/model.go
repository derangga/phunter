package tui

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"phunter/internal/process"
	"phunter/internal/theme"
)

const autoRefreshInterval = 5 * time.Second

// Mode represents the current UI mode.
type Mode int

const (
	ModeNormal Mode = iota
	ModeFilter
	ModeConfirmKill
)

// FilterField represents which filter input is active.
type FilterField int

const (
	FilterName FilterField = iota
	FilterPort
)

// SortKey represents which column to sort by.
type SortKey int

const (
	SortPID SortKey = iota
	SortProcess
	SortUser
	SortType
	SortPort
)

func (k SortKey) String() string {
	switch k {
	case SortPID:
		return "pid"
	case SortProcess:
		return "process"
	case SortUser:
		return "user"
	case SortType:
		return "type"
	case SortPort:
		return "port"
	}
	return "pid"
}

// Message types

type refreshMsg struct {
	processes []process.Process
	err       error
}

type tickMsg time.Time

type clearToastMsg struct{}

type killMsg struct {
	name string
	pid  int
	port string
	err  error
}

type killDoneMsg struct {
	pid int
}

type model struct {
	styles    Styles
	theme     theme.Theme
	selStyle  string // "bar" or "block"
	version   string

	allProcs  []process.Process // full snapshot
	viewProcs []process.Process // filtered + sorted
	cursor    int               // index into viewProcs
	offset    int               // first visible row

	mode        Mode
	filterField FilterField
	nameInput   textinput.Model
	portInput   textinput.Model

	sortKey SortKey
	sortAsc bool

	autoRefresh bool
	lastRefresh time.Time
	nextTickIn  int // seconds remaining

	killTarget int  // PID being confirmed
	toast      string
	toastUntil time.Time

	quitting bool
	width    int
	height   int
}

// New creates and returns the initial TUI model.
func New(cfg theme.Config, version string) model {
	styles := NewStyles(cfg.Theme)

	ni := textinput.New()
	ni.Placeholder = "filter by name or user..."
	ni.CharLimit = 64

	pi := textinput.New()
	pi.Placeholder = "filter by port..."
	pi.CharLimit = 16
	pi.Validate = func(s string) error {
		for _, r := range s {
			if r < '0' || r > '9' {
				return fmt.Errorf("digits only")
			}
		}
		return nil
	}

	return model{
		styles:    styles,
		theme:     cfg.Theme,
		selStyle:  cfg.SelectionStyle,
		version:   version,
		nameInput: ni,
		portInput: pi,
		sortKey:   SortPID,
		sortAsc:   true,
	}
}

func refreshCmd() tea.Cmd {
	return func() tea.Msg {
		procs, err := process.GetListeningPorts()
		return refreshMsg{processes: procs, err: err}
	}
}

func tickCmd() tea.Cmd {
	return tea.Tick(1*time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func clearToastCmd() tea.Cmd {
	return tea.Tick(3*time.Second, func(time.Time) tea.Msg {
		return clearToastMsg{}
	})
}

func (m model) Init() tea.Cmd {
	return tea.Batch(refreshCmd(), tickCmd())
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tickMsg:
		now := time.Now()
		elapsed := now.Sub(m.lastRefresh)
		if m.autoRefresh {
			remaining := int((autoRefreshInterval - elapsed).Seconds())
			if remaining < 0 {
				remaining = 0
			}
			m.nextTickIn = remaining
			if elapsed >= autoRefreshInterval {
				return m, tea.Batch(refreshCmd(), tickCmd())
			}
		}
		// Clear expired toast
		if m.toast != "" && now.After(m.toastUntil) {
			m.toast = ""
		}
		return m, tickCmd()

	case killMsg:
		if msg.err != nil {
			m.toast = fmt.Sprintf("error: %s", msg.err.Error())
			m.toastUntil = time.Now().Add(3 * time.Second)
			return m, tickCmd()
		}
		m.toast = fmt.Sprintf("killed %s (PID %d)", msg.name, msg.pid)
		m.toastUntil = time.Now().Add(3 * time.Second)
		return m, refreshCmd()

	case killDoneMsg:
		return m, nil

	case clearToastMsg:
		m.toast = ""
		return m, nil

	case refreshMsg:
		if msg.err != nil {
			m.toast = "refresh failed: " + msg.err.Error()
			m.toastUntil = time.Now().Add(3 * time.Second)
			return m, nil
		}
		m.allProcs = msg.processes
		m.lastRefresh = time.Now()
		m.nextTickIn = int(autoRefreshInterval.Seconds())
		m.recompute()
		return m, nil

	case tea.KeyMsg:
		return m.handleKey(msg)
	}

	return m, nil
}

func (m *model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch m.mode {
	case ModeFilter:
		return m.handleFilterKey(msg)
	case ModeConfirmKill:
		return m.handleConfirmKey(msg)
	default:
		return m.handleNormalKey(msg)
	}
}

func (m *model) handleNormalKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	key := msg.String()
	switch key {
	case "q", "ctrl+c":
		m.quitting = true
		return m, tea.Quit
	case "up":
		m.moveCursor(-1)
		return m, nil
	case "down":
		m.moveCursor(1)
		return m, nil
	case "g":
		m.cursor = 0
		m.offset = 0
		return m, nil
	case "G":
		if len(m.viewProcs) > 0 {
			m.cursor = len(m.viewProcs) - 1
			vh := m.viewHeight()
			if m.cursor >= vh {
				m.offset = m.cursor - vh + 1
			}
		}
		return m, nil
	case "/":
		m.mode = ModeFilter
		m.nameInput.Focus()
		m.filterField = FilterName
		return m, nil
	case "s":
		m.cycleSort()
		m.recompute()
		return m, nil
	case "enter", "k":
		if len(m.viewProcs) == 0 || m.cursor >= len(m.viewProcs) {
			return m, nil
		}
		m.killTarget = m.viewProcs[m.cursor].PID
		m.mode = ModeConfirmKill
		return m, nil
	case "r":
		m.toast = "refreshing..."
		m.toastUntil = time.Now().Add(2 * time.Second)
		return m, refreshCmd()
	case "a":
		m.autoRefresh = !m.autoRefresh
		if m.autoRefresh {
			m.lastRefresh = time.Now()
			m.nextTickIn = int(autoRefreshInterval.Seconds())
		}
		return m, nil
	}
	return m, nil
}

func (m *model) handleFilterKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	key := msg.String()
	switch key {
	case "esc":
		m.mode = ModeNormal
		m.nameInput.SetValue("")
		m.portInput.SetValue("")
		m.nameInput.Blur()
		m.portInput.Blur()
		m.recompute()
		return m, nil
	case "tab":
		if m.filterField == FilterName {
			m.filterField = FilterPort
			m.nameInput.Blur()
			m.portInput.Focus()
		} else {
			m.filterField = FilterName
			m.portInput.Blur()
			m.nameInput.Focus()
		}
		return m, nil
	case "enter":
		m.mode = ModeNormal
		m.nameInput.Blur()
		m.portInput.Blur()
		m.recompute()
		return m, nil
	}

	// Forward to active input
	var cmd tea.Cmd
	if m.filterField == FilterName {
		m.nameInput, cmd = m.nameInput.Update(msg)
	} else {
		m.portInput, cmd = m.portInput.Update(msg)
	}
	m.recompute()
	return m, cmd
}

func (m *model) handleConfirmKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "y":
		m.mode = ModeNormal
		// Find the process to kill
		for _, p := range m.allProcs {
			if p.PID == m.killTarget {
				return m, killCmd(p)
			}
		}
		return m, nil
	case "n", "esc":
		m.mode = ModeNormal
		return m, nil
	}
	return m, nil
}

func killCmd(p process.Process) tea.Cmd {
	return func() tea.Msg {
		err := process.Kill(p.PID)
		return killMsg{name: p.Name, pid: p.PID, port: p.Port, err: err}
	}
}

func (m *model) recompute() {
	m.viewProcs = applyFilterAndSort(
		m.allProcs,
		m.nameInput.Value(),
		m.portInput.Value(),
		m.sortKey,
		m.sortAsc,
	)
	// Clamp cursor
	if len(m.viewProcs) == 0 {
		m.cursor = 0
		m.offset = 0
	} else {
		if m.cursor >= len(m.viewProcs) {
			m.cursor = len(m.viewProcs) - 1
		}
		if m.cursor < 0 {
			m.cursor = 0
		}
		vh := m.viewHeight()
		if m.cursor < m.offset {
			m.offset = m.cursor
		}
		if m.cursor >= m.offset+vh {
			m.offset = m.cursor - vh + 1
		}
	}
}

func (m *model) moveCursor(delta int) {
	if len(m.viewProcs) == 0 {
		return
	}
	m.cursor += delta
	if m.cursor < 0 {
		m.cursor = 0
	}
	if m.cursor >= len(m.viewProcs) {
		m.cursor = len(m.viewProcs) - 1
	}
	vh := m.viewHeight()
	if m.cursor < m.offset {
		m.offset = m.cursor
	}
	if m.cursor >= m.offset+vh {
		m.offset = m.cursor - vh + 1
	}
}

func (m *model) cycleSort() {
	nextKey := SortKey((int(m.sortKey) + 1) % 5)
	if nextKey == SortPID && m.sortKey == SortPort {
		// Wrapped around — flip direction instead of advancing
		m.sortAsc = !m.sortAsc
	} else {
		m.sortKey = nextKey
		m.sortAsc = true
	}
}

func (m *model) viewHeight() int {
	// header(1) + header-border(1) + column-header(1) + status(1) + help(1) = 5
	// filter bar adds 3 when visible
	overhead := 5
	if m.mode == ModeFilter || m.nameInput.Value() != "" || m.portInput.Value() != "" {
		overhead += 3
	}
	vh := m.height - overhead
	if vh < 1 {
		vh = 1
	}
	return vh
}

// selectedProc returns the currently selected process, or nil if none.
func (m *model) selectedProc() *process.Process {
	if len(m.viewProcs) == 0 || m.cursor >= len(m.viewProcs) {
		return nil
	}
	p := m.viewProcs[m.cursor]
	return &p
}
