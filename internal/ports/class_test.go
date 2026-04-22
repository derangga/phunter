package ports

import "testing"

func TestClassify(t *testing.T) {
	tests := []struct {
		port int
		want Class
	}{
		{0, ClassAny},
		{22, ClassPrivileged},
		{80, ClassPrivileged},
		{443, ClassPrivileged},
		{2375, ClassPrivileged},
		{1, ClassPrivileged},
		{1023, ClassPrivileged},
		{3000, ClassDev},
		{3001, ClassDev},
		{5173, ClassDev},
		{5432, ClassDev},
		{6379, ClassDev},
		{8080, ClassDev},
		{8081, ClassDev},
		{8888, ClassDev},
		{11434, ClassDev},
		{1024, ClassRegistered},
		{4000, ClassRegistered},
		{49151, ClassRegistered},
		{49152, ClassEphemeral},
		{65535, ClassEphemeral},
	}
	for _, tt := range tests {
		got := Classify(tt.port)
		if got != tt.want {
			t.Errorf("Classify(%d) = %v, want %v", tt.port, got, tt.want)
		}
	}
}

func TestGlyph(t *testing.T) {
	tests := []struct {
		class Class
		want  string
	}{
		{ClassPrivileged, "◆"},
		{ClassDev, "●"},
		{ClassRegistered, "○"},
		{ClassEphemeral, "·"},
		{ClassAny, "·"},
	}
	for _, tt := range tests {
		got := tt.class.Glyph()
		if got != tt.want {
			t.Errorf("Class(%d).Glyph() = %q, want %q", tt.class, got, tt.want)
		}
	}
}
