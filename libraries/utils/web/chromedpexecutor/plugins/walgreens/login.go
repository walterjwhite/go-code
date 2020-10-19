package walgreens

import (
	"context"
	"github.com/chromedp/chromedp"
	"github.com/walterjwhite/go-application/libraries/utils/web/chromedpexecutor"
)

const (
	userIdField = "//*[@id=\"user-name\"]"
	userPasswordField = "//*[@id=\"user_password\"]"
	submitButton = "//*[@id=\"submit_btn\"]"

	securityQuestionRadio = "//*[@id=\"radio-security\"]"
	continueButton = "//*[@id=\"optionContinue\"]"

	securityQuestionField = "//*[@id=\"secQues\"]"
	securityAnswerField = "//*[@id=\"validate_security_answer\"]"
)

func (s *Session) Authenticate(ctx context.Context) {
	s.chromedpsession = chromedpexecutor.New(ctx)

	// no need to wait
	s.chromedpsession.Waiter = nil

	s.chromedpsession.Execute(chromedp.Navigate(uri))

	// username
	s.chromedpsession.Execute(
		//chromedp.SendKeys(userIdField, s.Credentials.Username),
		chromedp.WaitReady(userIdField),
		chromedp.SendKeys(userIdField, s.Credentials.Username),
		chromedp.Click(submitButton),
	)

	// password
	s.chromedpsession.Execute(
		chromedp.SendKeys(userPasswordField, s.Credentials.Password),
		chromedp.Click(submitButton),
	)

	// security question
	s.chromedpsession.Execute(
		chromedp.Click(securityQuestionRadio),
		chromedp.Click(continueButton),
	)

	// security question
	s.chromedpsession.Execute(
		chromedp.SendKeys(securityQuestionField, s.Credentials.SecretAnswer),
		chromedp.Click(securityAnswerField),
	)
}
