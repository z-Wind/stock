package server

type ParseFunc func(string) ([]byte, error)

type Request struct {
	Symbol    string
	ParseFunc ParseFunc
}

type _parseResult struct {
	symbol     string
	jsonResult []byte
	requests   []Request
	err        error
}

type Scheduler interface {
	Submit(Request)
	WorkerReady(chan Request)
	Run()
}
type FinalResult struct {
	Symbol   string
	JSONData []byte
	Error    error
}
