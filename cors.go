package socket

func (soc *socket) Cors(callback func(origin string) (err error)) {
	soc.eventsCors = append(soc.eventsCors, callback)
}
