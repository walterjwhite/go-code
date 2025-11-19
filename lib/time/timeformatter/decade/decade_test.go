package decade

import (
	"testing"
	"time"
)

func TestFormat(t *testing.T) {
	tests := []struct {
		year int
		want string
	}{
		{2021, "2021-2030"},
		{2020, "2011-2020"},
		{2025, "2021-2030"},
		{1999, "1991-2000"},
	}

	for _, tt := range tests {
		time := time.Date(tt.year, 1, 1, 0, 0, 0, 0, time.UTC)
		got := Format(time)
		if got != tt.want {
			t.Errorf("Format(%d) = %s; want %s", tt.year, got, tt.want)
		}
	}
}
