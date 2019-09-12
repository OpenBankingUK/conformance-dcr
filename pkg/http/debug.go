package http

import (
	"fmt"
	"net/http"
	"net/http/httputil"
)

func DebugRequest(req *http.Request) string {
	debug, err := httputil.DumpRequest(req, true)
	if err != nil {
		return fmt.Sprintf("cant debug request object: %s", err.Error())
	}
	return fmt.Sprintf("request:\n %s", string(debug))
}

func DebugResponse(r *http.Response) string {
	debug, err := httputil.DumpResponse(r, true)
	if err != nil {
		return fmt.Sprintf("cant debug response object: %s", err.Error())
	}
	return fmt.Sprintf("response:\n %s", string(debug))
}
