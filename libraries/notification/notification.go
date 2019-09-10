package notification

type NotificationType int

const (
	Error   NotificationType = -1
	Warning NotificationType = 0
	Info    NotificationType = 1
	Debug   NotificationType = 2
	Trace   NotificationType = 3
)

type Notification struct {
	Title       string
	Description string
	Icon        string
	AudioFile   string
	Type        NotificationType
}

type Notifier interface {
	Notify(notification Notification)
}

var notifier = New()

func OnCompletion() {
	if r := recover(); r != nil {
		// exception
		notifier.Notify(Notification{Title: "Error", Description: "Application execution completed abnormally.", Type: Error})
	} else {
		// normal
		notifier.Notify(Notification{Title: "Success", Description: "Application execution completed normal.", Type: Info})
	}
}
