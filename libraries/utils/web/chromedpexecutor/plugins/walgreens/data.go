package walgreens

import (
	"github.com/walterjwhite/go-application/libraries/utils/web/chromedpexecutor"
)

const (
	//uri = "https://photo.walgreens.com/store/welcome"
	//uri = "https://walgreens.com"
	uri = "https://www.walgreens.com/login.jsp"
)

type Credentials struct {
	Username     string
	Password     string
	SecretAnswer string
}

type Session struct {
	Credentials *Credentials

	chromedpsession *chromedpexecutor.ChromeDPSession
}
