package worker

type Worker interface {
	Name() string

	Work()

	ShortBreak()
	LongBreak()

	Lunch()

	Stop()
}

