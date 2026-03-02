package atsclient_test

import (
	"testing"

	atsclient "github.com/walterjwhite/go-code/lib/utils/ats-client"
	"github.com/walterjwhite/go-code/lib/utils/ats-client/eightfold"
	"github.com/walterjwhite/go-code/lib/utils/ats-client/greenhouse"
	"github.com/walterjwhite/go-code/lib/utils/ats-client/icims"
	"github.com/walterjwhite/go-code/lib/utils/ats-client/jobworx"
	"github.com/walterjwhite/go-code/lib/utils/ats-client/taleo"
	workdayats "github.com/walterjwhite/go-code/lib/utils/ats-client/workday"
)

func TestATSInterfaceImplementation(t *testing.T) {
	var _ atsclient.ATS = &workdayats.WorkdayATS{}
	var _ atsclient.ATS = &eightfold.EightfoldATS{}
	var _ atsclient.ATS = &greenhouse.GreenhouseATS{}
	var _ atsclient.ATS = &taleo.TaleoATS{}
	var _ atsclient.ATS = &jobworx.JobworxATS{}
	var _ atsclient.ATS = &icims.IcimsATS{}
}

func TestGetName(t *testing.T) {
	atsImplementations := map[string]atsclient.ATS{
		"workday":    &workdayats.WorkdayATS{},
		"eightfold":  &eightfold.EightfoldATS{},
		"greenhouse": &greenhouse.GreenhouseATS{},
		"taleo":      &taleo.TaleoATS{},
		"jobworx":    &jobworx.JobworxATS{},
		"icims":      &icims.IcimsATS{},
	}

	expectedNames := map[string]string{
		"workday":    "workday",
		"eightfold":  "eightfold",
		"greenhouse": "greenhouse",
		"taleo":      "taleo",
		"jobworx":    "jobworx",
		"icims":      "icims",
	}

	for key, ats := range atsImplementations {
		name := ats.GetName()
		expected := expectedNames[key]

		if name != expected {
			t.Errorf("Expected %s for %s, got %s", expected, key, name)
		}
	}
}
