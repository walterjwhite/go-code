package day

import (
	"testing"
	"time"
)

func TestFormat(t *testing.T) {
	tests := []struct {
		name string
		time time.Time
		want string
	}{
		{
			"default",
			time.Date(2024, 5, 15, 0, 0, 0, 0, time.UTC),
			"2021-2030/2024/05.May/15",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Format(tt.time); got != tt.want {
				t.Errorf("Format() = %v, want %v", got, tt.want)
			}
		})
	}
}
