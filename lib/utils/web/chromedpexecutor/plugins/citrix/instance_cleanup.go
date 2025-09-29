package citrix

func (i *Instance) onClose() {
	<-i.ctx.Done()
	i.cleanup()
}

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
