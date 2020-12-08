package target

import (
	"context"
	"github.com/chromedp/chromedp"
	"github.com/rs/zerolog/log"

	"github.com/walterjwhite/go/lib/utils/web/chromedpexecutor"
	"time"
)

const (
	// url = "https://login.target.com/gsp/static/v1/login/?client_id=ecom-web-1.0.0&ui_namespace=ui-default&back_button_action=browser&keep_me_signed_in=false&kmsi_default=false&actions=create_session_signin"
	url = "https://target.com"

	//*[@id="account"]/span[1]/span/div/svg/path
	// signInButton = "//*[@id=\"accountNav-signIn\"]/a/div"

	usernameField = "//*[@id=\"username\"]"
	passwordField = "//*[@id=\"password\"]"
)

func (s *Session) Login(ctx context.Context) {
	if s.chromedpsession != nil {
		s.Logout()
	}

	s.chromedpsession = chromedpexecutor.New(ctx)

	// no need to wait
	s.chromedpsession.Waiter = nil

	//defer s.Cancel()

	s.chromedpsession.Execute(
		chromedp.Navigate(url),
		chromedp.Click("//*[@id=\"account\"]/span[1]/span/div"),
		chromedp.WaitVisible("//*[@id=\"accountNav-signIn\"]/a/div"),
		chromedp.Click("//*[@id=\"accountNav-signIn\"]/a/div"),
		chromedp.SendKeys(usernameField, s.Credentials.Username),
		chromedp.SendKeys(passwordField, s.Credentials.Password),
		chromedp.Submit(passwordField),
	)

	s.chromedpsession.ExecuteTimeLimited(
		chromedpexecutor.TimeLimitedChromeAction{Action: chromedp.WaitVisible(logoutButton),
			Limit: 10 * time.Second, IsException: true, Message: "Login Failed"},
	)

	log.Info().Msgf("Successfully authenticated as: %v", s.Credentials.Username)
}
