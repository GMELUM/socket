package socket

func (soc *socket) Request(callback func(uri string) (err error)) {
	soc.eventsRequest = append(soc.eventsRequest, callback)
}
