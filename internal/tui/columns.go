package tui

import "github.com/charmbracelet/bubbles/table"

const (
	colPID  = 7
	colUser = 12
	colType = 7
	colAddr = 17
	colPort = 6
	// colName is computed dynamically from terminal width
)

func tableColumns(nameW int) []table.Column {
	return []table.Column{
		{Title: "PID", Width: colPID},
		{Title: "Process", Width: nameW},
		{Title: "User", Width: colUser},
		{Title: "Type", Width: colType},
		{Title: "Address", Width: colAddr},
		{Title: "Port", Width: colPort},
	}
}
