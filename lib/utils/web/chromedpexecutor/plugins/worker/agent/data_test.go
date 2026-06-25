package agent

import (
	"github.com/walterjwhite/go-code/lib/utils/ui/windows"
	"testing"
)

func TestConf_String(t *testing.T) {
	tests := []struct {
		name        string
		browserName string
		url         string
		expected    string
	}{
		{
			name:        "Standard Config",
			browserName: "chrome",
			url:         "https://google.com",
			expected:    "agent.Conf{chrome, https://google.com}",
		},
		{
			name:        "Empty Values",
			browserName: "",
			url:         "",
			expected:    "agent.Conf{, }",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Conf{
				BrowserName: tt.browserName,
				Url:         tt.url,
				WindowsConf: &windows.WindowsConf{},
				iteration:   5,
			}

			result := c.String()

			if result != tt.expected {
				t.Errorf("String() = %q, want %q", result, tt.expected)
			}
		})
	}
}
