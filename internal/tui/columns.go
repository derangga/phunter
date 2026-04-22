package tui

const (
	colGlyph = 2  // port class glyph + space
	colPID   = 7
	colUser  = 12
	colType  = 7
	colAddr  = 18
	colPort  = 8
	// colName is computed dynamically from terminal width
)

// columnTitles maps SortKey to column header text (ALL-CAPS).
var columnTitles = [5]string{"PID", "PROCESS", "USER", "TYPE", "PORT"}
