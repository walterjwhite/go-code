package worker

type Worker interface {
	String() string

	Work()

	ShortBreak()
	LongBreak()

	Lunch()

	Stop()
}

