package stocker

// ErrorNoSupport no support
type ErrorNoSupport struct {
	Message string
}

func (e ErrorNoSupport) Error() string {
	return e.Message
}
