package ports

import "github.com/charmbracelet/lipgloss"

// Class represents a port classification based on its number.
type Class int

const (
	ClassAny Class = iota
	ClassPrivileged
	ClassDev
	ClassRegistered
	ClassEphemeral
)

var devPorts = map[int]struct{}{
	3000: {}, 3001: {}, 5173: {}, 5432: {}, 6379: {},
	8080: {}, 8081: {}, 8888: {}, 11434: {},
}

// Classify returns the port class for a given port number.
func Classify(p int) Class {
	if p == 0 {
		return ClassAny
	}
	if p == 22 || p == 80 || p == 443 || p == 2375 {
		return ClassPrivileged
	}
	if p < 1024 {
		return ClassPrivileged
	}
	if _, ok := devPorts[p]; ok {
		return ClassDev
	}
	if p < 49152 {
		return ClassRegistered
	}
	return ClassEphemeral
}

// Glyph returns the single-character glyph for this port class.
func (c Class) Glyph() string {
	switch c {
	case ClassPrivileged:
		return "◆"
	case ClassDev:
		return "●"
	case ClassRegistered:
		return "○"
	case ClassAny, ClassEphemeral:
		return "·"
	}
	return "·"
}

// Color returns the lipgloss.Color for this port class using theme color strings.
func (c Class) Color(privileged, dev, registered, ephemeral, any string) lipgloss.Color {
	switch c {
	case ClassPrivileged:
		return lipgloss.Color(privileged)
	case ClassDev:
		return lipgloss.Color(dev)
	case ClassRegistered:
		return lipgloss.Color(registered)
	case ClassEphemeral:
		return lipgloss.Color(ephemeral)
	case ClassAny:
		return lipgloss.Color(any)
	}
	return lipgloss.Color(any)
}

// String returns the human-readable name.
func (c Class) String() string {
	switch c {
	case ClassPrivileged:
		return "privileged"
	case ClassDev:
		return "dev"
	case ClassRegistered:
		return "registered"
	case ClassEphemeral:
		return "ephemeral"
	case ClassAny:
		return "any"
	}
	return "any"
}
