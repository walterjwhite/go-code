package data

type Session struct {
	FeedId      string
	SessionPath string
	DataPath    string

	LatestReceivedRecord *Record
}
