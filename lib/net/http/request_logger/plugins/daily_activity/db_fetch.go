package daily_activity

import (
	"database/sql"

	"fmt"
)

func (c *Conf) fetchHTTPRequests() ([]string, []map[string]interface{}, error) {
	if c.db == nil {
		return nil, nil, fmt.Errorf("db is nil")
	}
	rows, err := c.db.Query("SELECT * FROM http_requests")
	if err != nil {
		return nil, nil, fmt.Errorf("query http_requests: %w", err)
	}
	defer close(rows)

	cols, err := rows.Columns()
	if err != nil {
		return nil, nil, fmt.Errorf("fetch columns: %w", err)
	}

	var records []map[string]interface{}
	vals := make([]interface{}, len(cols))
	ptrs := make([]interface{}, len(cols))
	for i := range vals {
		ptrs[i] = &vals[i]
	}

	for rows.Next() {
		if err := rows.Scan(ptrs...); err != nil {
			return nil, nil, fmt.Errorf("scan row: %w", err)
		}
		rec := make(map[string]interface{}, len(cols))
		for i, col := range cols {
			rec[col] = vals[i]
		}
		records = append(records, rec)
	}
	return cols, records, nil
}

func close(rows *sql.Rows) {
	_ = rows.Close()
}
