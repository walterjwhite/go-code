package walgreens

import (
	"context"
	"github.com/chromedp/chromedp"
	"github.com/walterjwhite/go-application/libraries/utils/web/chromedpexecutor"
)

func (s *Session) Authenticate(ctx context.Context) {
	s.chromedpsession = chromedpexecutor.New(ctx)

	// no need to wait
	s.chromedpsession.Waiter = nil

	s.chromedpsession.Execute(chromedp.Navigate(uri))

	/*
		s.chromedpsession.Execute(
			//chromedp.Click("//*[@id=\"signin-btn-header-2\"]"))
			chromedp.Click("/html/body/header/div/div[1]/div/section/div/nav[1]/div/div/div[1]/a/span[2]"),
			chromedp.Click("/html/body/header/div/div[1]/div/section/div/nav[1]/div/div/div[1]/ul/li[1]/a[1]/strong"))
	*/

	// username
	s.chromedpsession.Execute(
		chromedp.SendKeys("//*[@id=\"user-name\"]", s.Credentials.Username),
		chromedp.Click("//*[@id=\"submit_btn\"]"),
	)

	// password
	s.chromedpsession.Execute(
		chromedp.SendKeys("//*[@id=\"user_password\"]", s.Credentials.Password),
		chromedp.Click("//*[@id=\"submit_btn\"]"),
	)

	// security question
	s.chromedpsession.Execute(
		chromedp.Click("//*[@id=\"radio-security\"]"),
		chromedp.Click("//*[@id=\"optionContinue\"]"),
	)

	// security question
	s.chromedpsession.Execute(
		chromedp.SendKeys("//*[@id=\"secQues\"]", s.Credentials.SecretAnswer),
		chromedp.Click("//*[@id=\"validate_security_answer\"]"),
	)
}

func (s *Session) UploadPhotos() {
	s.chromedpsession.Execute(
		chromedp.Click("//*[@id=\"menu-photo\"]/a/span"),
		chromedp.Click("//*[@id=\"photoOrg-addPhotos-qmp-btn\"]"),
		chromedp.Click("//*[@id=\"po-fdropdown\"]/li[2]/a/span[1]/span[2]"),
	)

	// handle popup ...
}

func (s *Session) Logout() {
	s.chromedpsession.Execute(
		chromedp.Click("//*[@id=\"signOut\"]/strong"),
	)

	// handle popup ...
}
