package close_modal_window

type CloseModalWindow struct {
}

// check if the modal window exists
// TODO: invoke this before any operation because the other operation might be interrupted if we don't
func (c *CloseModalWindow) Exists() bool {
	return false
}

// click x to make modal window disappear
//     /html/body/div[1]/div/div/div[1]/span
