package socket

import (
	"github.com/gmelum/socket/entity/connect"
)

func (soc *socket) Reject(callback func(ch *connect.Connect)) {
	soc.eventReject = append(soc.eventReject, callback)
}
