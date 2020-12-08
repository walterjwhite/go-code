package spot

import (
	"github.com/walterjwhite/go/lib/external/spot/data"

	"github.com/walterjwhite/go/lib/application/property"
	"github.com/walterjwhite/go/lib/external/spot/action"
	"github.com/walterjwhite/go/lib/external/spot/action/daily_export"
	"github.com/walterjwhite/go/lib/external/spot/action/movement"
	"github.com/walterjwhite/go/lib/external/spot/action/new_record"
	"github.com/walterjwhite/go/lib/io/yaml"
)

func (c *Configuration) initActions() {
	c.Actions = []action.BackgroundAction{new_record.New(c.Session), daily_export.New(c.Session), movement.New(c.Session)}

	c.doInitActions(c.Actions)
}

func (c *Configuration) doInitActions(actions []action.BackgroundAction) {
	for _, action := range actions {
		yaml.Read(property.GetConfigurationFile(action, "spot", c.Session.FeedId, "actions"), action)
		action.Init(c.Session, c.ctx)
	}
}

func (c *Configuration) onNewRecord(old, new *data.Record) {
	for _, a := range c.Actions {
		recordAction, ok := a.(action.RecordAction)
		if ok {
			recordAction.OnNewRecord(old, new)
		}

	}
}
