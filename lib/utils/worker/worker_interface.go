package worker

type Worker interface {
	Work()

	ShortBreak()
	LongBreak()

	Lunch()

	Stop()
}
