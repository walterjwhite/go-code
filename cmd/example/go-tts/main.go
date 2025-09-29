package main

import (
	"github.com/go-tts/tts/pkg/audio"
	"github.com/go-tts/tts/pkg/speech"

	"github.com/walterjwhite/go-code/lib/application/logging"

	"strings"
)

func main() {
	s := audio.NewSpeaker()

	speechStream := speech.FromTextStream(strings.NewReader("Hello there!"), speech.LangEn)
	logging.Panic(s.Play(speechStream))
}
