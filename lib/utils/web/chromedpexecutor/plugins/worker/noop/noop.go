package noop

import (
	"context"
	"github.com/rs/zerolog/log"
)

func (i *State) Name() string {
	return "noop"
}

func (i *State) Init(ctx context.Context, headless bool, contextuals ...interface{}) error {
	return nil
}

func (i *State) Work(ctx context.Context, headless bool) {
	log.Debug().Msg("Work - NOOP")
}

func (i *State) Cleanup() {

}
