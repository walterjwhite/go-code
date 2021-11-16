package virginpulse

import (
	"github.com/walterjwhite/go/lib/utils/web/chromedpexecutor"
)

type Session struct {
	Uri string

	MenuXpath         string
	LogoffButtonXpath string

	UsernameXpath    string
	PasswordXpath    string
	LoginButtonXpath string

	Credentials *Credentials

	ByPassAuthentication bool

	Script []string

	ChromeDPSession *chromedpexecutor.ChromeDPSession
}

type Action struct {
	Name    string
	Actions []string
}

type Credentials struct {
	EmailAddress string
	Password     string
}

func (c *Credentials) SecretFields() []string {
	return []string{"EmailAddress", "Password"}
}
