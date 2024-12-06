package movement

import (
	ggeo "github.com/kellydunn/golang-geo"
	"github.com/walterjwhite/go-code/lib/external/geo"
	"github.com/walterjwhite/go-code/lib/external/spot/client"
	"github.com/walterjwhite/go-code/lib/external/spot/data"
	"strings"

	"time"
)

func (m *MovementConfiguration) hasMoved(old, new *data.Record) bool {
	if old == nil {
		return true
	}

	p1 := ggeo.NewPoint(old.Latitude, old.Longitude)
	p2 := ggeo.NewPoint(new.Latitude, new.Longitude)
	return geo.Distance(p1, p2) > m.MovementTolerance
}

func (m *MovementConfiguration) TimeSinceLastMovement() time.Duration {
	return time.Since(m.Session.LatestReceivedRecord.DateTime)
}

func isSuspend(r *data.Record) bool {
	return strings.Compare(string(client.OK), string(r.MessageType)) == 0
}
