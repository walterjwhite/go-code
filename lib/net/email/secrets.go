package email

func (c *EmailAccount) SecretFields() []string {
	return []string{"Username", "Password", "Domain", "EmailAddress.Address"}
}
