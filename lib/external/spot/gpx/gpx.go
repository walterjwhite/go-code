package gpx

import (
	"github.com/walterjwhite/go-code/lib/application/logging"
	"github.com/walterjwhite/go-code/lib/external/spot/data"
	ggpx "github.com/walterjwhite/go-garmin-gpx"
	"os"
	"path/filepath"
)


func Export(records []*data.Record, filename string) string {
	trackSegment := &ggpx.TrackSegment{}

	for _, r := range records {
		p := ggpx.TrackPoint{Latitude: ggpx.Latitude(r.Latitude), Longitude: ggpx.Longitude(r.Longitude)}
		trackSegment.TrackPoint = append(trackSegment.TrackPoint, p)
	}

	gpxData := &ggpx.GPX{Tracks: []ggpx.Track{{TrackSegments: []ggpx.TrackSegment{*trackSegment}}}}

	targetFilename := filename + ".gpx"

	logging.Panic(os.MkdirAll(filepath.Dir(targetFilename), os.ModePerm))

	logging.Panic(ggpx.Write(gpxData, targetFilename))

	return targetFilename
}
