package vanguard

func (c *Credentials) EncryptedFields() []string {
	return []string{"Username", "Password"}
}
