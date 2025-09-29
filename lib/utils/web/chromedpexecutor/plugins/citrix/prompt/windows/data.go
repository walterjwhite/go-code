package windows

type WindowsConf struct {
	Centered           bool
	AutomaticallyHides bool

	StartButtonHeight float64

	Version WindowsVersion
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
