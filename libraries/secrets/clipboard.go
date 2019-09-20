package secrets

import (
	"github.com/atotto/clipboard"
	"github.com/walterjwhite/go-application/libraries/logging"
)

func CopyToClipboard(secret string) {
	logging.Panic(clipboard.WriteAll(secret))
}
