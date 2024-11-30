package email

func (e *Provider) SecretFields() []string {
	return []string{"EmailSenderAccount.Username", "EmailSenderAccount.Password", "EmailSenderAccount.Domain", "EmailSenderAccount.EmailAddress.Name", "EmailSenderAccount.EmailAddress.Address"}
}
