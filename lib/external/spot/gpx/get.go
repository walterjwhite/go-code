package gpx

import (
	"github.com/walterjwhite/go/lib/application/logging"
	"github.com/walterjwhite/go/lib/external/spot/data"

	"bufio"
	"encoding/json"
	"os"
)

func get(filename string) []*data.Record {
	records := make([]*data.Record, 0)

	file, err := os.Open(filename)
	logging.Panic(err)

	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		r := &data.Record{}
		logging.Panic(json.Unmarshal([]byte(scanner.Text()), r))

		logging.Panic(scanner.Err())
		records = append(records, r)
	}

	return records
}
