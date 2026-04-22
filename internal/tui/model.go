package tui

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/table"
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

type batchKillMsg struct {
	killed int
	errors []string
}

type model struct {
	table    table.Model
	styles   Styles
	theme    theme.Theme
	selStyle string // "bar" or "block"
	version  string

	allProcs  []process.Process // full snapshot
	viewProcs []process.Process // filtered + sorted

	selectedPIDs map[int]bool // multi-select
	showHelp     bool         // floating help overlay

	mode        Mode
	filterField FilterField
	nameInput   textinput.Model
	portInput   textinput.Model

	sortKey SortKey
	sortAsc bool

	autoRefresh bool
	lastRefresh time.Time
	nextTickIn  int // seconds remaining

	killTarget int // PID being confirmed (-1 for batch)
	toast      string
	toastUntil time.Time

	quitting bool
	width    int
	height   int
}

// New creates and returns the initial TUI model.
func New(cfg theme.Config, version string) model {
	styles := NewStyles(cfg.Theme)

	t := table.New(
		table.WithColumns(tableColumns(20)),
		table.WithFocused(true),
		table.WithHeight(20),
	)
	t.SetStyles(styles.Table)

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
		table:        t,
		styles:       styles,
		theme:        cfg.Theme,
		selStyle:     cfg.SelectionStyle,
		version:      version,
		selectedPIDs: make(map[int]bool),
		nameInput:    ni,
		portInput:    pi,
		sortKey:      SortPID,
		sortAsc:      true,
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
		nameW := m.nameColWidth()
		m.table.SetColumns(m.sortedColumns(nameW))
		m.table.SetHeight(m.viewHeight())
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

	case batchKillMsg:
		m.selectedPIDs = make(map[int]bool)
		if len(msg.errors) > 0 {
			m.toast = fmt.Sprintf("killed %d, %d failed: %s", msg.killed, len(msg.errors), msg.errors[0])
		} else {
			m.toast = fmt.Sprintf("killed %d process(es)", msg.killed)
		}
		m.toastUntil = time.Now().Add(3 * time.Second)
		return m, refreshCmd()

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
	// Any key dismisses the help overlay first
	if m.showHelp {
		m.showHelp = false
		return m, nil
	}

	key := msg.String()
	switch key {
	case "q", "ctrl+c":
		m.quitting = true
		return m, tea.Quit
	case "?":
		m.showHelp = true
		return m, nil
	case " ":
		if len(m.viewProcs) == 0 {
			return m, nil
		}
		cursor := m.table.Cursor()
		if cursor < 0 || cursor >= len(m.viewProcs) {
			return m, nil
		}
		pid := m.viewProcs[cursor].PID
		if m.selectedPIDs[pid] {
			delete(m.selectedPIDs, pid)
		} else {
			m.selectedPIDs[pid] = true
		}
		m.table.SetRows(m.buildRows())
		return m, nil
	case "/":
		m.mode = ModeFilter
		m.nameInput.Focus()
		m.filterField = FilterName
		m.table.SetHeight(m.viewHeight())
		return m, nil
	case "s":
		m.cycleSort()
		m.recompute()
		return m, nil
	case "enter", "k":
		if len(m.viewProcs) == 0 {
			return m, nil
		}
		if len(m.selectedPIDs) > 0 {
			// Batch kill mode
			m.killTarget = -1
			m.mode = ModeConfirmKill
			return m, nil
		}
		cursor := m.table.Cursor()
		if cursor < 0 || cursor >= len(m.viewProcs) {
			return m, nil
		}
		m.killTarget = m.viewProcs[cursor].PID
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

	// Delegate navigation to the table
	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
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
		m.table.SetHeight(m.viewHeight())
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
		m.table.SetHeight(m.viewHeight())
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
		if len(m.selectedPIDs) > 0 {
			var targets []process.Process
			for _, p := range m.allProcs {
				if m.selectedPIDs[p.PID] {
					targets = append(targets, p)
				}
			}
			return m, batchKillCmd(targets)
		}
		// Single kill
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

func (m *model) recompute() {
	m.viewProcs = applyFilterAndSort(
		m.allProcs,
		m.nameInput.Value(),
		m.portInput.Value(),
		m.sortKey,
		m.sortAsc,
	)
	m.table.SetRows(m.buildRows())
	// Update column titles to reflect sort state
	nameW := m.nameColWidth()
	m.table.SetColumns(m.sortedColumns(nameW))
	m.table.SetHeight(m.viewHeight())
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
	// app header(1) + status(1) + help bar border(1) + help bar content(1) = 4
	// filter bar adds 3 when visible
	// The bubbles table handles its own header+border internally
	overhead := 4
	if m.mode == ModeFilter || m.nameInput.Value() != "" || m.portInput.Value() != "" {
		overhead += 3
	}
	vh := m.height - overhead
	if vh < 1 {
		vh = 1
	}
	return vh
}

func (m *model) nameColWidth() int {
	// bubbles table adds Padding(0,1) per cell = 2 chars per column
	// 7 columns × 2 = 14 chars of cell padding
	fixed := colGlyph + colPID + colUser + colType + colAddr + colPort + 14
	w := m.width - fixed
	if w < 12 {
		w = 12
	}
	return w
}

// sortedColumns returns tableColumns with the active sort column title annotated.
func (m *model) sortedColumns(nameW int) []table.Column {
	cols := tableColumns(nameW)
	// Column indices (0=glyph, 1=PID, 2=PROCESS, 3=USER, 4=TYPE, 5=ADDRESS, 6=PORT)
	sortColIdx := map[SortKey]int{
		SortPID:     1,
		SortProcess: 2,
		SortUser:    3,
		SortType:    4,
		SortPort:    6,
	}
	arrow := " ↑"
	if !m.sortAsc {
		arrow = " ↓"
	}
	if idx, ok := sortColIdx[m.sortKey]; ok {
		cols[idx].Title += arrow
	}
	return cols
}

// selectedProc returns the currently selected process, or nil if none.
func (m *model) selectedProc() *process.Process {
	cursor := m.table.Cursor()
	if len(m.viewProcs) == 0 || cursor < 0 || cursor >= len(m.viewProcs) {
		return nil
	}
	p := m.viewProcs[cursor]
	return &p
}
