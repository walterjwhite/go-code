package main

import (
	//	"github.com/chromedp/cdproto/input"
	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/kb"
	"github.com/walterjwhite/go-application/libraries/application"
	"github.com/walterjwhite/go-application/libraries/chromedpexecutor"
	"time"
)

var (
	s *chromedpexecutor.ChromeDPSession
)

func init() {
	application.Configure()

	s = chromedpexecutor.LaunchRemoteBrowser(application.Context)
}

func main() {
	// do not wait
	s.Waiter = nil

	// do some stuff
	s.Execute(chromedp.Navigate("https://duckduckgo.com"))

	// this tabs off of the search bar

	// Alt key is innocous and triggers a response in Windows applications (so we can see it is doing something)
	//s.Execute(chromedp.KeyEvent(kb.Alt))
	//s.Execute(chromedp.KeyEvent("test"))

	//moveMouse()
	tabThroughElements()
}

func tabThroughElements() {
	for {
		s.Execute(chromedp.KeyEvent(kb.Tab))
		time.Sleep(400 * time.Millisecond)
	}
}

/*
const (
	radius = 500
)

func moveMouse() {
	// indefinite
	for {
		for x := 0; x <= radius; x++ {
			move(x, 0)
		}
		for x := radius; x >= 0; x-- {
			move(x, 0)
		}
	}
}

func move(x, y int) {
	s.Execute(chromedp.MouseEvent(input.MouseMoved, float64(x), float64(y)))
	time.Sleep(10 * time.Millisecond)
}
*/
