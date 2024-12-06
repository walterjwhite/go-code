package database

import (
	"github.com/walterjwhite/go-code/lib/application/logging"
)

type CountQuery struct {
	Query       Query
	RecordCount int
}

func (q *CountQuery) Count() {
	q.Query.connect()

	rows, err := q.Query.Database.Query(q.Query.QueryString, q.Query.Parameters)
	logging.Panic(err)

	defer rows.Close()

	if rows.Next() {
		var count int
		err = rows.Scan(&count)
		logging.Panic(err)

		q.RecordCount = count
	}
}
