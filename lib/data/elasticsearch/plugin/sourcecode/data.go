package sourcecode

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/walterjwhite/go-code/lib/application/logging"
	"strconv"
	"strings"
)

/*
 * generic so that it can be used for:
 * 1. code search
 * 2. reference search
 */
type Line struct {
	Path string
	Name string

	Number int

	// project, branch
	Tags []string

	Contents string
}

func (l *Line) DocumentId() string {
	sha256Sum := sha256.Sum256([]byte(fmt.Sprintf("%v.%v.%v.%v", l.Path, l.Name, l.Number, strings.Join(l.Tags, "."))))
	return hex.EncodeToString(sha256Sum[:])
}

/*
func (l *Line) Mapping() string {

}
*/

func (l *Line) Equals(record []string) bool {
	number, err := strconv.Atoi(record[2])
	logging.Panic(err)

	//tags := strings.Split(record[3], ",")

	return l.Path == record[0] && l.Name == record[1] && l.Number == number /*&& l.Tags == tags*/
}
