package typename

import (
	"fmt"
	"strings"
)

func Get(data any) string {
	return strings.ReplaceAll(fmt.Sprintf("%T", data), "*", "")
}
