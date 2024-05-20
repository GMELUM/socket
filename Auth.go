package socket

func (soc *socket) Auth(callback func(url []byte) (err error)) {
	soc.eventsAuth = append(soc.eventsAuth, callback)
}
