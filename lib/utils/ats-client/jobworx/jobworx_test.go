package jobworx

import (
	"testing"

	atsclient "github.com/walterjwhite/go-code/lib/utils/ats-client"
)

func TestJobworxGetName(t *testing.T) {
	jobworx := &JobworxATS{}
	name := jobworx.GetName()

	if name != "jobworx" {
		t.Errorf("Expected name to be 'jobworx', got '%s'", name)
	}
}

func TestJobworxImplementsInterface(t *testing.T) {
	var _ atsclient.ATS = &JobworxATS{}
}

