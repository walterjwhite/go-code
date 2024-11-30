package gateway

import (
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/application/logging"
)

func (s *Session) Validate() {
	if len(s.Credentials.Domain) == 0 {
		logging.Panic(errors.New("Domain is required"))
	}
	if len(s.Credentials.Username) == 0 {
		logging.Panic(errors.New("Username is required"))
	}
	if len(s.Credentials.Password) == 0 {
		logging.Panic(errors.New("Password is required"))
	}
	if len(s.Credentials.Pin) == 0 {
		logging.Panic(errors.New("Pin is required"))
	}

	log.Info().Msg("Validated session configuration")
}

func validateToken(token string) {
	if len(token) != 6 {
		logging.Panic(fmt.Errorf("Please enter the 6-digit token: %v", token))
	}
}
