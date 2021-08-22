package patientnotebook

func (s *Session) SecretFields() []string {
	return []string{"Credentials.Username", "Credentials.Password"}
}
