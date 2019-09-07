package database

import (
	"github.com/jmoiron/sqlx"
	"log"
)

type Query struct {
	Query      string
	Parameters []string

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
		err := q.Database.Select(dest, q.Query, q.Parameters)
		if err != nil {
			log.Fatalf("Error querying: %v / %v\n", q.Query, err)
		}
	} else {
		err := q.Database.Select(dest, q.Query)
		if err != nil {
			log.Fatalf("Error querying: %v / %v\n", q.Query, err)
		}
	}
}
