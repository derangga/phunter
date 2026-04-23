package tui

import (
	"testing"

	"phunter/internal/process"
)

func TestApplyFilterAndSort(t *testing.T) {
	procs := []process.Process{
		{PID: 100, Name: "nginx", User: "root", Type: "IPv4", Address: "0.0.0.0", Port: "80"},
		{PID: 200, Name: "node", User: "dev", Type: "IPv4", Address: "127.0.0.1", Port: "3000"},
		{PID: 300, Name: "postgres", User: "postgres", Type: "IPv4", Address: "127.0.0.1", Port: "5432"},
		{PID: 400, Name: "sshd", User: "root", Type: "IPv4", Address: "0.0.0.0", Port: "22"},
		{PID: 150, Name: "vite", User: "dev", Type: "IPv4", Address: "127.0.0.1", Port: "5173"},
	}

	t.Run("no filter", func(t *testing.T) {
		result := applyFilterAndSort(procs, "", "", SortPID, true)
		if len(result) != 5 {
			t.Fatalf("expected 5, got %d", len(result))
		}
		if result[0].PID != 100 || result[4].PID != 400 {
			t.Errorf("sort by PID asc failed: first=%d last=%d", result[0].PID, result[4].PID)
		}
	})

	t.Run("name filter", func(t *testing.T) {
		result := applyFilterAndSort(procs, "node", "", SortPID, true)
		if len(result) != 1 || result[0].Name != "node" {
			t.Errorf("name filter failed: got %d results", len(result))
		}
	})

	t.Run("name filter case insensitive", func(t *testing.T) {
		result := applyFilterAndSort(procs, "NGINX", "", SortPID, true)
		if len(result) != 1 {
			t.Errorf("case insensitive filter expected 1, got %d", len(result))
		}
	})

	t.Run("port filter", func(t *testing.T) {
		result := applyFilterAndSort(procs, "", "54", SortPID, true)
		if len(result) != 1 || result[0].Port != "5432" {
			t.Errorf("port filter expected postgres, got %d results", len(result))
		}
	})

	t.Run("combined filter (AND)", func(t *testing.T) {
		result := applyFilterAndSort(procs, "dev", "30", SortPID, true)
		if len(result) != 1 || result[0].Name != "node" {
			t.Errorf("combined filter expected node, got %d results", len(result))
		}
	})

	t.Run("sort by process asc", func(t *testing.T) {
		result := applyFilterAndSort(procs, "", "", SortProcess, true)
		if result[0].Name != "nginx" || result[4].Name != "vite" {
			t.Errorf("sort by process asc: first=%s last=%s", result[0].Name, result[4].Name)
		}
	})

	t.Run("sort by port desc", func(t *testing.T) {
		result := applyFilterAndSort(procs, "", "", SortPort, false)
		if result[0].Port != "5432" {
			t.Errorf("sort by port desc: first port=%s, expected 5432", result[0].Port)
		}
	})

	t.Run("sort by PID desc", func(t *testing.T) {
		result := applyFilterAndSort(procs, "", "", SortPID, false)
		if result[0].PID != 400 {
			t.Errorf("sort by PID desc: first=%d, expected 400", result[0].PID)
		}
	})
}

func TestHighlightMatch(t *testing.T) {
	tests := []struct {
		text, substr         string
		before, match, after string
	}{
		{"Hello World", "world", "Hello ", "World", ""},
		{"nginx", "gin", "n", "gin", "x"},
		{"nginx", "NGINX", "", "nginx", ""},
		{"nginx", "", "nginx", "", ""},
		{"nginx", "xyz", "nginx", "", ""},
	}
	for _, tt := range tests {
		b, m, a := highlightMatch(tt.text, tt.substr)
		if b != tt.before || m != tt.match || a != tt.after {
			t.Errorf("highlightMatch(%q, %q) = (%q, %q, %q), want (%q, %q, %q)",
				tt.text, tt.substr, b, m, a, tt.before, tt.match, tt.after)
		}
	}
}
