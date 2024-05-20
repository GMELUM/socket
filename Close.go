package socket

func (soc *socket) Close() {
	soc.listener.Close()
}
