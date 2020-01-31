package database

import (
	databasel "github.com/walterjwhite/go-application/libraries/database"
	"github.com/walterjwhite/go-application/libraries/logging"
)

type Source struct {
	// driver
	// connection information
	// query
	Query databasel.Query
}

func (s *Source) Write(channel chan interface{}) {
	rows, err := s.Query.Queryx()
	logging.Panic(err)

	defer s.Query.Database.DB.Close()
	for rows.Next() {
		var d interface{}
		logging.Panic(rows.StructScan(&d))

		// push to channel
		channel <- d
	}
}
