package idle

import (
	"github.com/walterjwhite/go-application/libraries/application/logging"
	"syscall"
	"time"
	"unsafe"
)

var (
	user32           = syscall.MustLoadDLL("user32.dll")
	kernel32         = syscall.MustLoadDLL("kernel32.dll")
	getLastInputInfo = user32.MustFindProc("GetLastInputInfo")
	getTickCount     = kernel32.MustFindProc("GetTickCount")
	lastInputInfo    struct {
		cbSize uint32
		dwTime uint32
	}
)

func IdleTime() time.Duration {
	lastInputInfo.cbSize = uint32(unsafe.Sizeof(lastInputInfo))
	currentTickCount, _, _ := getTickCount.Call()

	rl, _, err := getLastInputInfo(uintptr(unsafe.Pointer(&lastInputInfo)))
	if rl == 0 {
		logging.Panic(err)
	}

	return time.Duration(uint32(currentTickCount) - lastInputInfo.dwTime*time.Millisecond)
}
