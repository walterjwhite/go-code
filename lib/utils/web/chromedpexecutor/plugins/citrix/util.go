package citrix

import (
	"strings"
)

func Get(input string) []int {
	token := make([]int, 6)

	for i, char := range input {
		value := int(char - '0')

		token[i] = value
	}

	return token
}

func (s *Session) trim(token string) string {
	s.Credentials.Username = strings.TrimSpace(s.Credentials.Username)
	s.Credentials.Domain = strings.TrimSpace(s.Credentials.Domain)
	s.Credentials.Password = strings.TrimSpace(s.Credentials.Password)

	s.Credentials.Pin = strings.TrimSpace(s.Credentials.Pin)

	return strings.TrimSpace(token)
}
