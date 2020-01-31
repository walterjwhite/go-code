package stage3

// sets up chroot environment so that the rest of the modules may work
func (m *Stage3) Run() error {
	m.setupRoot()
	m.configureMakeConf()
}

func (m *Stage3) setupRoot() {
	if !m.isChrootSetup() {

	}

	if !m.isStage3Valid() {
		m.downloadStage3()

		if !m.isStage3Valid() {
			return errors.New("Stage 3 is invalid")
		}
	}
}

func (m *Stage3) isChrootSetup() bool {

}
