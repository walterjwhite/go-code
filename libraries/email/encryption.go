package email

func (c *EmailSenderAccount) EncryptedFields() []string {
	return []string{"Username", "Password", "Domain"}
}
