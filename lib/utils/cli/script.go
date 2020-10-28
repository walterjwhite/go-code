package cli

import (
	"context"
	"github.com/rs/zerolog/log"
	"path/filepath"
	"strconv"
	"time"
)

func Run(ctx context.Context, scriptFile *ScriptFile, parentDirectory string) {
	for i, command := range scriptFile.Commands {
		log.Debug().Msgf("running command(%v)", i)

		command.LogDirectory = filepath.Join(parentDirectory, strconv.Itoa(i))
		command.Execute(ctx)

		if scriptFile.DelayBetweenCommands != nil && i < (len(scriptFile.Commands)-1) {
			log.Debug().Msgf("delaying between commands(%v)", *scriptFile.DelayBetweenCommands)
			time.Sleep(*scriptFile.DelayBetweenCommands)
		}
	}
}
