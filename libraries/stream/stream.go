package stream

type Source interface {
	// read data from underlying source and push to channel
	Write(channel chan interface{})
}

type Sink interface {
	Read(channel chan interface{})
}

// create a channel that connects the Sink to the Source
func Pipe(source Source, sink Sink) {
	channel := make(chan interface{})

	go source.Write(channel)
	go sink.Read(channel)
}
