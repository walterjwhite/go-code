package craigslist

type CraigslistPost struct {
	Region string

	Seller   OwnerType
	Category Category

	Title       string
	Description string
	Price       string

	City       string
	PostalCode string

	/*
	Make      string
	Model     string
	Size      string
	Condition ConditionType
	*/
	Script []string
	// images to attach
	Media []string

	EmailAddress     string
	PhoneNumber      string
	ReceiveTexts     bool
	ReceiveCalls     bool
	PhoneContactName string
}
