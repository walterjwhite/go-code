package main

import (
	"github.com/go-vgo/robotgo"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-application/libraries/application"
	"time"
)

func init() {
	application.Configure()
}

func main() {
	for {
		//robotgo.MoveMouseSmooth(100, 200, 1.0, 100.0)
		log.Debug().Msgf("Moving mouse to: %v, %v", 1700, 200)
		robotgo.MoveMouse(1700, 200)

		time.Sleep(1 * time.Second)

		log.Debug().Msgf("Moving mouse to: %v, %v", 1800, 200)
		robotgo.MoveMouse(1800, 200)

		time.Sleep(1 * time.Second)

		//robotgo.MoveMouseSmooth(100, 200, 1.0, 100.0)
	}

	/*
		//robotgo.TypeStr("Hello World")
		robotgo.TypeString("Hello World")
		robotgo.Sleep(1)

		//robotgo.KeyTap("enter")
		robotgo.KeyTap("i", "alt", "command")
	*/

	/*
		x, y := robotgo.GetMousePos()
		log.Printf("@ %d,%d", x, y)

		color := robotgo.GetPixelColor(x, y)
		log.Printf("color %v", color)
	*/

	/*
		bitmap := robotgo.CaptureScreen(10, 20, 30, 40)
		// free the bitmap
		defer robotgo.FreeBitmap(bitmap)

		fmt.Println("...", bitmap)

		// find this location
		fx, fy := robotgo.FindBitmap(bitmap)
		fmt.Println("Find ...", fx, fy)

		robotgo.SaveBitmap(bitmap, "test.png")
	*/

	// program blocks until you press the keys in this order ...
	/*
		ok := robotgo.AddEvents("q", "ctrl", "shift")
		if ok {
			fmt.Println("add events ...")
		}
	*/

	/*
		keve := robotgo.AddEvent("k")
		if keve {
			fmt.Println("you pressed...", "k")
		}

		mleft := robotgo.AddEvent("mleft")
		if mleft {
			fmt.Println("You pressed ...", "mouse left button")
		}
	*/

	/*
		fpid, err := robotgo.FindIds("chrom")
		if err == nil {
			fmt.Println("pids...", fpid)

			if len(fpid) > 0 {
				fmt.Println("active pid:0", robotgo.ActivePID(fpid[0]))

				//robotgo.Kill(4290)
			}

			abool := robotgo.ShowAlert("test", "robotgo")
			if abool == 0 {
				fmt.Println("ok@@@", "ok")
			}

			title := robotgo.GetTitle()
			fmt.Println("title@@@ ", title)
		}
		* */
}
