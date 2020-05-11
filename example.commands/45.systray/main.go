package main

import (
	"github.com/getlantern/systray"
)

func main() {
	// Should be called at the very beginning of main().
	systray.Run(onReady, onExit)
}

func onReady() {
	//	systray.SetIcon(icon.Data)
	systray.SetTitle("Awesome App")
	systray.SetTooltip("Pretty awesome超级棒")
	//	mQuit := systray.AddMenuItem("Quit", "Quit the whole app")

	// Sets the icon of a menu item. Only available on Mac.
	//	mQuit.SetIcon(icon.Data)
}

func onExit() {
	// clean up here
}
