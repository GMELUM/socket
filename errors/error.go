package errors

import "github.com/gobwas/ws"

func Error(code int) error {
	return ws.RejectConnectionError(
		ws.RejectionStatus(403),
	)
}
