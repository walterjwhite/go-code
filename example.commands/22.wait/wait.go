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

func init() {
	application.Configure()
}

func main() {
	interval := 1 * time.Second
	limit := 5 * time.Second
	d := &data{index: 0}

	wait.Wait(application.Context, &interval, &limit, d.f)
}

func (d *data) f() bool {
	defer d.increment()

	log.Info().Msgf("checking: %v", d.index)
	return d.index > 2
}

func (d *data) increment() {
	log.Info().Msgf("incrementing: %v", d.index)
	d.index++
}
