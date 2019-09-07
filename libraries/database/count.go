package database

import (
	"log"
)

type CountQuery struct {
	Query       Query
	RecordCount int
}

func (q *CountQuery) Count() {
	q.Query.connect()

	// automatically inject parameters
	rows, err := q.Query.Database.Query(q.Query.Query, q.Query.Parameters)
	if err != nil {
		log.Fatalf("Error querying: %v / %v\n", q.Query.Query, err)
	}

	defer rows.Close()

	// iterate over each row
	if rows.Next() {
		var count int
		err = rows.Scan(&count)
		if err != nil {
			log.Printf("Error scanning rows: %v\n", err)
		}

		q.RecordCount = count
	}
}
