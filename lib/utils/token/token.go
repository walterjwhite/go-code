package token

type TokenProvider interface {
	Get() string
}
