package gateway

import (
	"github.com/chromedp/chromedp"
)

// TODO: move this into a go-plugin / module
func (s *Session) launchDesktop() {
	s.chromedpsession.Execute(
		chromedp.Click("//*[@id=\"desktopsBtn\"]/div"),
		chromedp.Click("//*[@id=\"home-screen\"]/div[2]/section[4]/div[4]/ul/li/a[2]"),
		chromedp.Click("//*[@id=\"appInfoOpenBtn\"]"),
	)
}

// TODO: move this into a go-plugin / module
func (s *Session) launchRemoteDesktop() {
	s.chromedpsession.Execute(
		chromedp.Click("//*[@id=\"allAppsBtn\"]"),
		chromedp.Click("//*[@id=\"home-screen\"]/div[2]/section[4]/div[4]/ul/li[15]/a[1]/div[3]/p[1]"),
		chromedp.Click("//*[@id=\"appInfoOpenBtn\"]"),
	)

	// wait 5 seconds (minimum)
	// 10 seconds (maximum)

	// send hostname
	//s.chromedpsession.Execute(

	/*
			https://groups.google.com/forum/#!topic/chrome-debugging-protocol/DQxlrBNSC9w
			options = {
		  "type": "keyDown",
		  "key": "Tab"
		}


		chrome.debugger.sendCommand({tabId:debuggee}, "Input.dispatchKeyEvent",options,function(b){
		 options.type = "keyUp"
		 chrome.debugger.sendCommand({tabId:debuggee}, "Input.dispatchKeyEvent",options)
		})
	*/

	// hit connect
	// send password
	// hit enter

}
