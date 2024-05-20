package socket

func (soc *socket) Cors(callback func(origin []byte) (err error)) {
	soc.eventsCors = append(soc.eventsCors, callback)
}
