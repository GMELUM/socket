package socket

import (
	"github.com/gmelum/socket/entity/connect"
)

func (soc *socket) Disconnect(callback func(ch *connect.Connect)) {
	soc.eventsDisconnect = append(soc.eventsDisconnect, callback)
}
