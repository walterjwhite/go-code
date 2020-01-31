package module

type Module interface {
	Document() error
	Run() error
}
