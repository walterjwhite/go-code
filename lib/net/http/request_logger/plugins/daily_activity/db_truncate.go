package daily_activity

import (
	"fmt"
)

func (c *Conf) truncateHTTPRequests() error {
	if c.db == nil {
		return fmt.Errorf("db is nil")
	}
	if _, err := c.db.Exec("DELETE FROM http_requests"); err != nil {
		return fmt.Errorf("truncate http_requests: %w", err)
	}
	return nil
}
