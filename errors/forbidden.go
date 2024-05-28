package errors

func Forbidden() error {
	return Error(403)
}
