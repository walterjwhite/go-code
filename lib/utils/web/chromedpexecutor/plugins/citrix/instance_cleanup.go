package citrix

func (i *Instance) cleanup() {
	if i.Worker != nil {
		i.Worker.Cleanup()
	}

	select {
	case <-i.ctx.Done():
		return
	default:
		i.cancel()
	}
}
