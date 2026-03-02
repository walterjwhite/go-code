package taleo

import (
	"testing"

	atsclient "github.com/walterjwhite/go-code/lib/utils/ats-client"
)

func TestTaleoGetName(t *testing.T) {
	taleo := &TaleoATS{}
	name := taleo.GetName()

	if name != "taleo" {
		t.Errorf("Expected name to be 'taleo', got '%s'", name)
	}
}

func TestTaleoImplementsInterface(t *testing.T) {
	var _ atsclient.ATS = (*TaleoATS)(nil)
}
