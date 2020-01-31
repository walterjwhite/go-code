package approval

import (
	"encoding/csv"
	"github.com/walterjwhite/go-application/libraries/logging"
	//"github.com/walterjwhite/go-application/libraries/timeformatter/timestamp"
	"os"
)

// todo; follow the activity pattern
// everything gets written out into a project
// git support is provided
// select either one/many records / file
type csvApprovalWriter struct {
	// "/tmp/approval"
	OutputFilename string
}

func (w *csvApprovalWriter) Write(action *Action) {
	f, err := os.OpenFile(w.OutputFilename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	logging.Panic(err)

	defer f.Close()

	// Other []string
	writer := csv.NewWriter(f)
	logging.Panic(writer.Write([]string{action.Request.Id, timestamp.Get(), action.Action, action.Comments, action.Username}))
	writer.Flush()

	// user: TODO: get this from the request headers
}
