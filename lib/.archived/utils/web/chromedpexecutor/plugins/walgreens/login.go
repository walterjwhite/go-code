package walgreens

import (
	"context"
	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/kb"
	"github.com/walterjwhite/go/lib/utils/web/chromedpexecutor"
	// "time"
)

// const (
// 	securityQuestionRadio = "//*[@id=\"radio-security\"]"
// 	continueButton        = "//*[@id=\"optionContinue\"]"

// 	securityQuestionField          = "//*[@id=\"secQues\"]"
// 	validateSecurityQuestionButton = "//*[@id=\"validate_security_answer\"]"
// )

func (s *Session) Login(ctx context.Context) {
	s.chromedpsession = chromedpexecutor.New(ctx)

	// no need to wait
	s.chromedpsession.Waiter = nil

	s.chromedpsession.Execute(chromedp.Navigate(url))

	s.chromedpsession.Execute(
		chromedp.KeyEvent(s.Credentials.Username),
		chromedp.KeyEvent(kb.Tab),
		chromedp.KeyEvent(kb.Tab),
		chromedp.KeyEvent(s.Credentials.Password),
		chromedp.KeyEvent(kb.Enter),
	)

	//s.answerSecurityQuestion(ctx)

	// this isn't working much as logging in by finding elements also does not work
	// s.chromedpsession.ExecuteTimeLimited(
	// 	chromedpexecutor.TimeLimitedChromeAction{Action: chromedp.WaitVisible(loggedInField),
	// 		Limit: 10 * time.Second, IsException: true, Message: "Login Failed"},
	// )
}

// func (s *Session) answerSecurityQuestion(ctx context.Context) {
// 	// select security question
// 	s.chromedpsession.Execute(
// 		chromedp.Click(securityQuestionRadio),
// 		chromedp.Click(continueButton),
// 	)

// 	// security question
// 	s.chromedpsession.Execute(
// 		chromedp.SendKeys(securityQuestionField, s.Credentials.SecretAnswer),
// 		chromedp.Click(validateSecurityQuestionButton),
// 	)
// }
