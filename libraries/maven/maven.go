package maven

const COMMAND = "mvn"

const QUIET_LOGS = "-Dorg.slf4j.simpleLogger.defaultLogLevel=WARN"
const QUIET_TRANSFERS = "-ntp"

func GetCommandLine(arguments []string, debug *bool) (string, []string) {
	if !*debug {
		arguments = append(arguments, QUIET_LOGS)
		arguments = append(arguments, QUIET_TRANSFERS)
	}

	return COMMAND, arguments
}
