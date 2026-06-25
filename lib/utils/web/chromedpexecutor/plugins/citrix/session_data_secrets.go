package citrix

func (s *Session) SecretFields() []string {
	return []string{"Credentials.Username", "Credentials.Password", "Credentials.Domain", "Credentials.Pin"}
}
