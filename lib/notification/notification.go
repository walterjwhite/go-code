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

var NotifierInstance = New()

func OnCompletion() {
	if r := recover(); r != nil {
		NotifierInstance.Notify(Notification{Title: "Error", Description: "Application execution completed abnormally.", Type: Error})
	} else {
		NotifierInstance.Notify(Notification{Title: "Success", Description: "Application execution completed normally.", Type: Info})
	}
}
