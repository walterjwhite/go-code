package timestamp

import (
	"testing"
	"time"
)

func TestFormat(t *testing.T) {
	testTime := time.Date(2024, 5, 15, 10, 30, 0, 0, time.UTC)

	Default = &Configuration{Template: "%d.%d.%d.%d.%d.%d.%d"}
	want := "2024.5.5.15.10.30.0"
	got := Format(testTime)
	if got != want {
		t.Errorf("Default Format() = %s; want %s", got, want)
	}

	UseNested()
	want = "2024/05.May/15/10.30.0"
	got = Format(testTime)
	if got != want {
		t.Errorf("Nested Format() = %s; want %s", got, want)
	}

	Default = &Configuration{Template: "%d.%d.%d.%d.%d.%d.%d"}
}
