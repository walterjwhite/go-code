package application

import (
	"bytes"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"strings"
	"testing"
)

func TestIsConfigured(t *testing.T) {
	oldVersion, oldBuild := ApplicationVersion, BuildDate
	defer func() {
		ApplicationVersion = oldVersion
		BuildDate = oldBuild
	}()

	ApplicationVersion = ""
	BuildDate = ""
	if isConfigured() {
		t.Fatalf("expected not configured when version and build date empty")
	}

	ApplicationVersion = "v1"
	BuildDate = "2025-12-10"
	if !isConfigured() {
		t.Fatalf("expected configured when version and build date set")
	}
}

func TestGetApplicationId(t *testing.T) {
	oldName, oldVersion, oldScm := ApplicationName, ApplicationVersion, SCMId
	defer func() {
		ApplicationName = oldName
		ApplicationVersion = oldVersion
		SCMId = oldScm
	}()

	ApplicationName = "myapp"
	ApplicationVersion = "v2"
	SCMId = "abc123"

	got := GetApplicationId()
	want := "myapp.v2.abc123"
	if got != want {
		t.Fatalf("GetApplicationId() = %q; want %q", got, want)
	}
}

func TestLogIdentifier_NotConfigured(t *testing.T) {
	oldVersion, oldBuild := ApplicationVersion, BuildDate
	defer func() {
		ApplicationVersion = oldVersion
		BuildDate = oldBuild
	}()

	ApplicationVersion = ""
	BuildDate = ""

	var buf bytes.Buffer
	log.Logger = zerolog.New(&buf)

	logIdentifier()

	if !strings.Contains(buf.String(), "Application was not built properly") {
		t.Fatalf("expected log message not found")
	}
}
