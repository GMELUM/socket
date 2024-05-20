package socket

import (
	"github.com/gmelum/socket/entity/connect"
)

func (soc *socket) Reject(callback func(ch *connect.Connect)) {
	soc.eventsReject = append(soc.eventsReject, callback)
}
