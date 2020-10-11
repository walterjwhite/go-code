package walgreens

func (c *Credentials) EncryptedFields() []string {
	return []string{"Username", "Password", "SecretAnswer"}
}
