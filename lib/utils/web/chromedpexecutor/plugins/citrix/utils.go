package citrix

import (
	"strings"
)

func (s *Session) trim(token string) string {
	s.Credentials.Username = strings.TrimSpace(s.Credentials.Username)
	s.Credentials.Domain = strings.TrimSpace(s.Credentials.Domain)
	s.Credentials.Password = strings.TrimSpace(s.Credentials.Password)

	s.Credentials.Pin = strings.TrimSpace(s.Credentials.Pin)

	return strings.TrimSpace(token)
}
