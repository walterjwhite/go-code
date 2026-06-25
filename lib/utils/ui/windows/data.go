package windows

import (
	"fmt"
	"image"

	"github.com/walterjwhite/go-code/lib/utils/ui"
)

type WindowsConf struct {
	Centered           bool
	AutomaticallyHides bool

	StartButtonHeight float64

	Version WindowsVersion

	Controller ui.Controller
}

type Application struct {
	Name                 string
	TaskBarIconUpImage   image.Image
	TaskBarIconDownImage image.Image
	TitleBarIconImage    image.Image

	WindowsConf *WindowsConf
}

func (c *WindowsConf) String() string {
	return fmt.Sprintf("windowsConf{%v, %v, %f, %d}", c.Centered, c.AutomaticallyHides, c.StartButtonHeight, c.Version)
}

func (a *Application) String() string {
	return fmt.Sprintf("application{%s}", a.Name)
}

type WindowsVersion int

const (
	Windows10 WindowsVersion = iota
	Windows11
)

func Windows10Default() *WindowsConf {
	return &WindowsConf{Centered: false, AutomaticallyHides: false, StartButtonHeight: 48, Version: Windows10}
}

func Windows11Default() *WindowsConf {
	return &WindowsConf{Centered: true, AutomaticallyHides: false, StartButtonHeight: 48, Version: Windows11}
}
