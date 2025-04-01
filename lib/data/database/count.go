package database

import (
	"database/sql"
	"github.com/walterjwhite/go-code/lib/application/logging"
)

type CountQuery struct {
	Query       Query
	RecordCount int
}

func (q *CountQuery) Count() {
	q.Query.connect()

	rows, err := q.Query.Database.Query(q.Query.QueryString, q.Query.Parameters)
	defer q.cleanup(rows)

	logging.Panic(err)

	if rows.Next() {
		var count int
		err = rows.Scan(&count)
		logging.Panic(err)

		q.RecordCount = count
	}
}

func (q *CountQuery) cleanup(rows *sql.Rows) {
	if rows == nil {
		return
	}

	logging.Panic(rows.Close())
}
