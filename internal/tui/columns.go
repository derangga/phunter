package tui

import "github.com/charmbracelet/bubbles/table"

const (
	colGlyph = 2 // port class glyph + space
	colPID   = 7
	colUser  = 12
	colType  = 7
	colAddr  = 18
	colPort  = 8
	// colName is computed dynamically from terminal width
)

func tableColumns(nameW int) []table.Column {
	return []table.Column{
		{Title: " ", Width: colGlyph},
		{Title: "PID", Width: colPID},
		{Title: "PROCESS", Width: nameW},
		{Title: "USER", Width: colUser},
		{Title: "TYPE", Width: colType},
		{Title: "ADDRESS", Width: colAddr},
		{Title: "PORT", Width: colPort},
	}
}
