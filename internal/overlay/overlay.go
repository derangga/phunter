package overlay

import (
	"strings"
	"unicode/utf8"

	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-runewidth"
)

// Place stamps fg on top of bg at visual position (x, y).
func Place(x, y int, fg, bg string) string {
	fgLines := strings.Split(fg, "\n")
	bgLines := strings.Split(bg, "\n")
	for i, fgLine := range fgLines {
		row := y + i
		if row < 0 || row >= len(bgLines) {
			continue
		}
		bgLine := bgLines[row]
		bgW := lipgloss.Width(bgLine)
		fgW := lipgloss.Width(fgLine)

		left := takeColumns(bgLine, x)
		leftW := lipgloss.Width(left)
		if leftW < x {
			left += strings.Repeat(" ", x-leftW)
		}

		right := ""
		if x+fgW < bgW {
			right = "\x1b[0m" + skipColumns(bgLine, x+fgW)
		}

		bgLines[row] = left + "\x1b[0m" + fgLine + right
	}
	return strings.Join(bgLines, "\n")
}

// takeColumns returns the first n visual columns of s, preserving ANSI sequences.
func takeColumns(s string, n int) string {
	col, i := 0, 0
	for i < len(s) {
		if s[i] == '\x1b' {
			j := i + 1
			if j < len(s) && s[j] == '[' {
				j++
				for j < len(s) && !(s[j] >= '@' && s[j] <= '~') {
					j++
				}
				if j < len(s) {
					j++
				}
			}
			i = j
			continue
		}
		r, size := utf8.DecodeRuneInString(s[i:])
		w := runewidth.RuneWidth(r)
		if col+w > n {
			break
		}
		col += w
		i += size
	}
	return s[:i]
}

// skipColumns returns the portion of s starting at visual column n.
func skipColumns(s string, n int) string {
	col, i := 0, 0
	for i < len(s) && col < n {
		if s[i] == '\x1b' {
			j := i + 1
			if j < len(s) && s[j] == '[' {
				j++
				for j < len(s) && !(s[j] >= '@' && s[j] <= '~') {
					j++
				}
				if j < len(s) {
					j++
				}
			}
			i = j
			continue
		}
		r, size := utf8.DecodeRuneInString(s[i:])
		col += runewidth.RuneWidth(r)
		i += size
	}
	return s[i:]
}
