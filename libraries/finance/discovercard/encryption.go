package discovercard

func (c *WebCredentials) EncryptedFields() []string {
	return []string{"Username", "Password"}
}
