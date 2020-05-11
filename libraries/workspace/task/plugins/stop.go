package plugins

type PreStop interface {
	PreStop()
}

type PostStop interface {
	PostStop()
}
