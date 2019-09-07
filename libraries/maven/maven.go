package maven

const COMMAND = "mvn"

const QUIET_LOGS = "-Dorg.slf4j.simpleLogger.defaultLogLevel=WARN"

func GetCommandLine(arguments []string, debug *bool) (string, []string) {
	if !*debug {
		arguments = append(arguments, QUIET_LOGS)
	}

	return COMMAND, arguments
}
