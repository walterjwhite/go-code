package codesearch

import (
	"fmt"
	"github.com/walterjwhite/go-application/libraries/logging"
	"io"
)

type SearchOutputProcessor interface {
	OnMatch(match *Match)
}

type DefaultSearchOutputProcessor struct {
	Output io.Writer
}

func (p *DefaultSearchOutputProcessor) OnMatch(match *Match) {
	_, err := p.Output.Write([]byte(fmt.Sprintf("%v:%v:%v", match.Filename, match.LineNumber, string(match.Matched))))
	logging.Panic(err)
}
