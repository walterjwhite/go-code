package bnymellon401k

func (s *Session) SecretFields() []string {
	return []string{"Credentials.Username", "Credentials.Password"}
}
