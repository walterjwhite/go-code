package learning

func (s *Session) SecretFields() []string {
	return []string{"EmailAddress", "Password"}
}
