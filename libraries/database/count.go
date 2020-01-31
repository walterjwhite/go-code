package database

import (
	"github.com/walterjwhite/go-application/libraries/logging"
)

type CountQuery struct {
	Query       Query
	RecordCount int
}

func (q *CountQuery) Count() {
	q.Query.connect()

	// automatically inject parameters
	rows, err := q.Query.Database.Query(q.Query.QueryString, q.Query.Parameters)
	logging.Panic(err)

	defer rows.Close()

	// iterate over each row
	if rows.Next() {
		var count int
		err = rows.Scan(&count)
		logging.Panic(err)

		q.RecordCount = count
	}
}
