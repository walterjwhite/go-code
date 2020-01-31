package gmail

import (
	"context"
	"fmt"
	"github.com/chromedp/chromedp"
	"github.com/walterjwhite/go-application/libraries/chromedpexecutor"
	"github.com/walterjwhite/go-application/libraries/token"
	"strconv"
)

func (a *Account) doCreate(ctx context.Context, tokenProvider token.TokenProvider) {
	a.session = chromedpexecutor.New(ctx)

	a.session.Execute(chromedp.Navigate(gmailBaseUrl))
	a.session.Execute(a.doCompleteForm()...)

	a.session.Execute(a.doSelectVerificationMethod()...)
	a.session.Execute(a.doVerify(tokenProvider)...)
	a.session.Execute(a.doPersonalInformation()...)
	a.session.Execute(a.doEula()...)
}

func (a *Account) doCompleteForm() []chromedp.Action {
	var actions []chromedp.Action

	actions = append(actions, chromedp.SendKeys("//*[@id=\"firstName\"]", a.FirstName))
	actions = append(actions, chromedp.SendKeys("//*[@id=\"lastName\"]", a.LastName))
	actions = append(actions, chromedp.SendKeys("//*[@id=\"username\"]", a.Username))

	// password
	actions = append(actions, chromedp.SendKeys("//*[@id=\"passwd\"]/div[1]/div/div[1]/input", a.Password))
	actions = append(actions, chromedp.SendKeys("//*[@id=\"confirm-passwd\"]/div[1]/div/div[1]/input", a.Password))

	// next
	actions = append(actions, chromedp.Click("//*[@id=\"accountDetailsNext\"]/div[2]"))

	// phoneNumber
	actions = append(actions, chromedp.SendKeys("//*[@id=\"phoneNumberId\"]", a.PhonePreference.PhoneNumber))

	// next
	actions = append(actions, chromedp.Click("//*[@id=\"gradsIdvPhoneNext\"]/div[2]"))

	return actions
}

func (a *Account) doSelectVerificationMethod() []chromedp.Action {
	if !a.PhonePreference.Call {
		return nil
	}

	// call instead

	var actions []chromedp.Action
	actions = append(actions, chromedp.Click("//*[@id=\"view_container\"]/form/div[2]/div/div[4]/div[1]/div[2]/button"))

	return actions
}

func (a *Account) doVerify(tokenProvider token.TokenProvider) []chromedp.Action {
	// read token
	// TODO: use a token-service here (can be implemented by reading from stdin, file, etc.)
	code := tokenProvider.Get()

	var actions []chromedp.Action

	actions = append(actions, chromedp.SendKeys("//*[@id=\"code\"]", code))

	// TODO: potential exception here
	// verified this is on the screen
	actions = append(actions, chromedp.Click("//*[@id=\"gradsIdvVerifyNext\"]/span/span"))

	// TODO: potential exception here
	// this is not on the screen
	// finish
	actions = append(actions, chromedp.Click("//*[@id=\"view_container\"]/form/div[2]/div/div[2]/div/div[1]/div/div[1]/input"))

	return actions
}

func (a *Account) doPersonalInformation() []chromedp.Action {
	var actions []chromedp.Action

	actions = append(actions, chromedp.SendKeys("//*[@id=\"month\"]", strconv.Itoa(a.BirthDate.Month)))
	actions = append(actions, chromedp.SendKeys("//*[@id=\"day\"]", strconv.Itoa(a.BirthDate.Day)))
	actions = append(actions, chromedp.SendKeys("//*[@id=\"year\"]", strconv.Itoa(a.BirthDate.Year)))

	actions = append(actions, chromedp.Click(fmt.Sprintf("//*[@id=\"gender\"]/option[%v]", a.Gender)))

	actions = append(actions, chromedp.Click("//*[@id=\"personalDetailsNext\"]/span/span"))

	return actions
}

func (a *Account) doEula() []chromedp.Action {
	var actions []chromedp.Action

	actions = append(actions, chromedp.Click("//*[@id=\"personalDetailsNext\"]/span/span"))

	if !a.PhonePreference.PhoneUsage {
		// skip
		actions = append(actions, chromedp.Click("//*[@id=\"view_container\"]/form/div[2]/div/div[2]/div[1]/div[2]/button"))
	} else {
		// yes im in
		actions = append(actions, chromedp.Click("//*[@id=\"phoneUsageNext\"]/span/span"))
	}

	// scroll to bottom ...
	// click until it disappears (3x+)
	actions = append(actions, chromedp.Click("//*[@id=\"view_container\"]/form/div[2]/div/div/div/div[1]/div/div/span/span/svg"))
	actions = append(actions, chromedp.Click("//*[@id=\"view_container\"]/form/div[2]/div/div/div/div[1]/div/div/span/span/svg"))
	actions = append(actions, chromedp.Click("//*[@id=\"view_container\"]/form/div[2]/div/div/div/div[1]/div/div/span/span/svg"))

	// confirm EULA
	actions = append(actions, chromedp.Click("//*[@id=\"termsofserviceNext\"]/span/span"))

	return actions
}
