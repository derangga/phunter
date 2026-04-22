package tui

import (
	"sort"
	"strconv"
	"strings"

	"phunter/internal/process"
)

// applyFilterAndSort filters and sorts the process list.
func applyFilterAndSort(all []process.Process, nameFilter, portFilter string, key SortKey, asc bool) []process.Process {
	nameFilter = strings.ToLower(strings.TrimSpace(nameFilter))
	portFilter = strings.TrimSpace(portFilter)

	out := make([]process.Process, 0, len(all))
	for _, p := range all {
		if nameFilter != "" {
			if !strings.Contains(strings.ToLower(p.Name), nameFilter) &&
				!strings.Contains(strings.ToLower(p.User), nameFilter) {
				continue
			}
		}
		if portFilter != "" {
			if !strings.Contains(p.Port, portFilter) {
				continue
			}
		}
		out = append(out, p)
	}

	sort.SliceStable(out, func(i, j int) bool {
		if asc {
			return lessProc(out[i], out[j], key)
		}
		return lessProc(out[j], out[i], key)
	})
	return out
}

func lessProc(a, b process.Process, key SortKey) bool {
	switch key {
	case SortPID:
		return a.PID < b.PID
	case SortProcess:
		return strings.ToLower(a.Name) < strings.ToLower(b.Name)
	case SortUser:
		return strings.ToLower(a.User) < strings.ToLower(b.User)
	case SortType:
		return a.Type < b.Type
	case SortPort:
		ap := portToInt(a.Port)
		bp := portToInt(b.Port)
		return ap < bp
	}
	return a.PID < b.PID
}

func portToInt(s string) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		return -1 // "*" or non-numeric → sort first
	}
	return n
}

// highlightMatch splits text around a case-insensitive substring match,
// returning (before, match, after). If no match, match and after are empty.
func highlightMatch(text, substr string) (string, string, string) {
	if substr == "" {
		return text, "", ""
	}
	lower := strings.ToLower(text)
	idx := strings.Index(lower, strings.ToLower(substr))
	if idx < 0 {
		return text, "", ""
	}
	return text[:idx], text[idx : idx+len(substr)], text[idx+len(substr):]
}
