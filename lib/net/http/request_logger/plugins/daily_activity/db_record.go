package daily_activity

import (
	"fmt"

	"time"
)

func convertRecordToStrings(cols []string, rec map[string]interface{}) []string {
	out := make([]string, len(cols))
	for i, c := range cols {
		v := rec[c]
		if v == nil {
			out[i] = ""
			continue
		}
		switch t := v.(type) {
		case []byte:
			out[i] = string(t)
		case string:
			out[i] = t
		case time.Time:
			out[i] = t.Format(time.RFC3339)
		default:
			out[i] = fmt.Sprintf("%v", t)
		}
	}
	return out
}
