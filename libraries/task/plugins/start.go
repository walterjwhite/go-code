package plugins

type PreStart interface {
	PreStart()
}

type PostStart interface {
	PostStart()
}
