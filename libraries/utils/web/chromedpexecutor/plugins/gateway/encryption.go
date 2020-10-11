package gateway

func (c *Credentials) EncryptedFields() []string {
	return []string{"Username", "Password", "Domain", "Pin"}
}
