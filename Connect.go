package socket

import (
	"github.com/gmelum/socket/entity/connect"
)

func (soc *socket) Connect(callback func(ch *connect.Connect)) {
	soc.eventsConnect = append(soc.eventsConnect, callback)
}
