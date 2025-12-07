package request_logger

import (
	"github.com/walterjwhite/go-code/lib/application/logging"
)

const (
	INSERT_STATEMENT = `
	INSERT INTO http_requests (ts, ip, method, request_uri, user_agent, status)
	VALUES (:ts, :ip, :method, :request_uri, :user_agent, :status)
	`
)

func (c *Conf) initDB() {
	statements := c.Provider.InitSQL()
	for i := range statements {
		_, _ = c.DB.Exec(statements[i])
	}
}

func (c *Conf) initRequestLogChannel() chan<- RequestLog {
	ch := make(chan RequestLog, c.BufferSize)

	go func() {
		for requestLogRecord := range ch {
			_, err := c.DB.NamedExec(INSERT_STATEMENT, requestLogRecord)
			logging.Warn(err, "initRequestLogChannel - insert")
		}
	}()

	return ch
}
