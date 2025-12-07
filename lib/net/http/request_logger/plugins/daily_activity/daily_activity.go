package daily_activity

import (
	"github.com/jmoiron/sqlx"

	"github.com/walterjwhite/go-code/lib/application/logging"

	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/net/email"
)

type Conf struct {
	db           *sqlx.DB
	emailAccount *email.EmailAccount
}

func With(db *sqlx.DB, emailAccount *email.EmailAccount) *Conf {
	return &Conf{db: db, emailAccount: emailAccount}
}

func (c *Conf) ExportAndEmailRequests() {
	cols, rows, err := c.fetchHTTPRequests()
	if err != nil {
		logging.Warn(err, "ExportAndEmailRequests.fetchHTTPRequests")
		return
	}
	if len(rows) == 0 {
		log.Info().Msg("daily_activity: no rows to export")
		return
	}

	body, err := generateCSVBody(cols, rows)
	if err != nil {
		logging.Warn(err, "ExportAndEmailRequests.generateCSVBody")
		return
	}

	logging.Warn(c.sendEmail(body), "ExportAndEmailRequests.sendEmail")
	logging.Warn(c.truncateHTTPRequests(), "ExportAndEmailRequests.truncateHTTPRequests")

	log.Info().Msg("daily_activity: exported and truncated http_requests")
}
