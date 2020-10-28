package database

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/walterjwhite/go/lib/application/logging"
)

type Query struct {
	QueryString string
	Parameters  []string

	Database                *sqlx.DB
	ConnectionConfiguration ConnectionConfiguration
}

func (q *Query) connect() {
	if q.Database == nil {
		q.Database = q.ConnectionConfiguration.Connect()
	}
}

func (q *Query) Select(dest interface{}) {
	q.connect()

	defer q.Database.DB.Close()

	if q.Parameters != nil {
		logging.Panic(q.Database.Select(dest, q.QueryString, q.Parameters))
	} else {
		logging.Panic(q.Database.Select(dest, q.QueryString))
	}
}

func (q *Query) Query() (*sql.Rows, error) {
	q.connect()

	if q.Parameters != nil {
		return q.Database.Query(q.QueryString, q.Parameters)
	}

	return q.Database.Query(q.QueryString)
}

func (q *Query) Queryx() (*sqlx.Rows, error) {
	q.connect()

	if q.Parameters != nil {
		return q.Database.Queryx(q.QueryString, q.Parameters)
	}

	return q.Database.Queryx(q.QueryString)
}
