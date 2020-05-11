package secrets

type Configuration struct {
	Keys             []string
	SupportedRegions []string
}

type Profile struct {
	Configurations []*Configuration
}
