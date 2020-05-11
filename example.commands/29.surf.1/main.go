package main

import (
	//"fmt"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-application/libraries/application"
	"github.com/walterjwhite/go-application/libraries/logging"
	//"gopkg.in/headzoo/surf.v1"
	//"github.com/headzoo/surf.v1.0.0"
	//"github.com/headzoo/surf/browser.v1.0.0"
	"github.com/headzoo/surf"
	"github.com/headzoo/surf/browser"
	//"gopkg.in/headzoo/surf.v1/browser"
	//	"gopkg.in/headzoo/surf/browser.v1"
)

func init() {
	application.Configure()
}

func main() {
	bow := surf.NewBrowser()
	logging.Panic(bow.Open("https://discovercard.com"))

	// Outputs: "The Go Programming Language"
	log.Info().Msg(bow.Title())

	login(bow)
	navigateToCreditCardActivity(bow)
}

func login(bow *browser.Browser) {
	fm, err := bow.Form( /*"login-form-content"*/ /* "form"*/ "#login-form-content")
	logging.Panic(err)

	logging.Panic(fm.Input( /*"userid-content"*/ "userID", ""))
	logging.Panic(fm.Input( /*"password-content"*/ "password", ""))
	logging.Panic(fm.Submit())

	log.Info().Msg(bow.Title())
	log.Info().Msg(bow.Body())
}

func navigateToCreditCardActivity(bow *browser.Browser) {
	//logging.Panic(bow.Click("*[@id=\"main-content\"]/div[3]/div[3]/div/div[2]/strong"))

	//logging.Panic(bow.Click("#main-content div:nth-child(3) div:nth-child(3) div div:nth-child(1) a div p:nth-child(1)"))
	logging.Panic(bow.Click("#main-content > div:nth-of-type(3) > div:nth-of-type(3) > div > div:nth-of-type(1) > a > div > p:nth-of-type(1)"))
	log.Info().Msg(bow.Title())
	log.Info().Msg(bow.Body())

	logging.Panic(bow.Click("*[@id=\"main-content-rwd\"]/div[4]/div[1]/div[1]/div[1]/p[1]/span"))
	log.Info().Msg(bow.Title())
	log.Info().Msg(bow.Body())
}

//*[@id="sso-portal-moneymarket"]/div/p[1]
//*[@id="sso-portal-moneymarket"]/div/p[1]
//*[@id="sso-portal-moneymarket"]/div/p[1]
