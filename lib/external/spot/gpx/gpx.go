package gpx

import (
	"github.com/walterjwhite/go-code/lib/application/logging"
	"github.com/walterjwhite/go-code/lib/external/spot/data"
	ggpx "github.com/walterjwhite/go-garmin-gpx"
	"os"
	"path/filepath"
)

// NOTES:
// 1. this library writes to out/filename.gpx
// 2. timestamp information isn't written out, it might be useful
// 3. source file *MUST* have 1 line / record and not have the record over multiple lines

// returns the path to the newly written file (matching how the gpx library writes the file)
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
