package response

import "encoding/json"

const RequestIDHeader = "request_id"

type HTTPResponse struct {
	Code      int    `json:"code"`
	Message   string `json:"message"`
	RequestID string `json:"request_id"`
	TraceID   string `json:"trace_id"`
	Data      any    `json:"data"`
}

func NewHTTPResponse(bs []byte) *HTTPResponse {
	var res = new(HTTPResponse)
	_ = json.Unmarshal(bs, res)
	return res
}
