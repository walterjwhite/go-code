package worker

import (
  "github.com/rs/zerolog/log"
  "time"
)

func (c *Conf) manageBreaks() {
  if c == nil || c.OnDuration == nil {
    log.Warn().Msg("disabling taking breaks")
    return
  }

  if c.LunchStartTime.Till()+*c.LunchDuration < 0 {
    c.hadLunch = true
  }

  for {
    time.Sleep(*c.OnDuration)

    if !c.hadLunch && c.LunchStartTime.Till() < 0 {
      c.hadLunch = true
      c.doBreak("time for lunch: %v", c.LunchDuration)
      c.cycle = 0

      continue
    }

    if c.cycle%4 == 0 {
      c.doBreak("taking break: %v", c.BreakDuration)
    } else {
      c.doBreak("taking off: %v", c.OffDuration)
    }

    c.cycle++
  }
}

func (c *Conf) doBreak(message string, duration *time.Duration) {
  log.Info().Msgf(message, *duration)

  c.breakChannel <- duration
  time.Sleep(*duration)
}
