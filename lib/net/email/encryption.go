package email

func (c *EmailSenderAccount) SecretFields() []string {
	return []string{"Username", "Password", "Domain"}
}
