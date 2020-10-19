package gateway

func (c *Credentials) SecretFields() []string {
	return []string{"Username", "Password", "Domain", "Pin"}
}
