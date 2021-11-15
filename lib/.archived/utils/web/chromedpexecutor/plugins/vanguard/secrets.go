package vanguard

func (s *Session) SecretFields() []string {
	return []string{"Credentials.Username", "Credentials.Password"}
}
