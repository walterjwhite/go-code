func (i *Instance) waitForDesktop() {
	screenCheckDuration := 5 * time.Second
	if !i.isScreenLocked() {
		log.Warn().Msg("screen is not locked (windows icon is present)")
		return
	}

	for {
		if i.isWaitingForTermsAcceptance() {
			i.handlePrompt()
		} else if i.isLoggingIn() {
			log.Info().Msg("waiting for system to login")
			for {
				time.Sleep(screenCheckDuration)
				if !i.isLoggingIn() {
					log.Info().Msg("logged in")
					return
				}
			}
		} else {
			log.Warn().Msg("neither waiting for terms acceptance or logging in")
		}

		time.Sleep(screenCheckDuration)
	}
}
