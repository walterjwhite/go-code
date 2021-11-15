package walgreens

func (s *Session) SecretFields() []string {
	return []string{"Credentials.Username", "Credentials.Password", "Credentials.SecretAnswer"}
}
