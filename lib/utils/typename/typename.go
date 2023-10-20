package typename

import (
	"fmt"
	"strings"
)

func Get(data interface{}) string {
	return strings.ReplaceAll(fmt.Sprintf("%T", data), "*", "")
}
