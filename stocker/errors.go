package stocker

// ErrorNoSupport no support
type ErrorNoSupport struct {
	Message string
}

func (e ErrorNoSupport) Error() string {
	return e.Message
}

// ErrorNotFound no found
type ErrorNoFound struct {
	Message string
}

func (e ErrorNoFound) Error() string {
	return e.Message
}

// ErrorFatal fatal error
type ErrorFatal struct {
	Message string
}

func (e ErrorFatal) Error() string {
	return e.Message
}
