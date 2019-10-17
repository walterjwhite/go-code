package main

import (
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-application/libraries/application"
	"github.com/walterjwhite/go-application/libraries/wait"
	"time"
)

type data struct {
	index int
}

func main() {
	ctx := application.Configure()

	interval := 1 * time.Second
	limit := 2 * time.Second
	d := &data{index: 0}

	wait.Wait(ctx, interval, limit, d.f)
}

func (d *data) f() bool {
	defer d.increment()

	log.Info().Msgf("checking: %v", d.index)
	return d.index > 5
}

func (d *data) increment() {
	log.Info().Msgf("incrementing: %v", d.index)
	d.index++
}
