package maven

const (
	command = "mvn"

	quietLogs      = "-Dorg.slf4j.simpleLogger.defaultLogLevel=WARN"
	quietTransfers = "-ntp"
)

func GetCommandLine(arguments []string, debug *bool) (string, []string) {
	if !*debug {
		arguments = append(arguments, quietLogs)
		arguments = append(arguments, quietTransfers)
	}

	return command, arguments
}
