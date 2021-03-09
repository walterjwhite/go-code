package main

import (
	"context"
	"github.com/walterjwhite/go/lib/application"

	"github.com/walterjwhite/go/lib/utils/web/chromedpexecutor/plugins/discovercard"
)

var (
	session = &discovercard.Session{}
)

func init() {
	application.ConfigureWithProperties(session)
}

func main() {
	defer application.OnEnd()

	session.Login(context.Background())

	//time.Sleep(10 * time.Minute)

	// get balance
	//*[@id="main-content"]/div[3]/div[3]/div/div[1]/a/div/p[1]
	//*[@id="main-content"]/div[3]/div[3]/div/div[1]/strong

	//*[@id="sso-portal-moneymarket"]/div/p[1]
	//*[@id="sso-portal-moneymarket"]/div/p[2]
	//*[@id="main-content"]/div[4]/div[3]/div[1]/div[1]/strong

	//*[@id="sso-portal-moneymarket"]/div
	//*[@id="sso-portal-moneymarket"]/div/p[1]
	//*[@id="sso-portal-moneymarket"]/div/p[2]
	//*[@id="main-content"]/div[4]/div[3]/div[2]/div[1]/strong

	//*[@id="sso-portal-moneymarket"]/div
}
