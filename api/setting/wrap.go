package setting

import "net/http"

const (
	userAgent = `Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:62.0) Gecko/20100101 Firefox/62.0`
)

// WrapRequest add header
func WrapRequest(req *http.Request) *http.Request {
	req.Header.Add("User-Agent", userAgent)
	return req
}
