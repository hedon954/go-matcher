package response

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HTTPResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func NewHTTPResponse(bs []byte) *HTTPResponse {
	var res = new(HTTPResponse)
	_ = json.Unmarshal(bs, res)
	return res
}

func GinParamError(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, HTTPResponse{
		Code:    http.StatusBadRequest,
		Message: err.Error(),
		Data:    nil,
	})
}

func GinError(c *gin.Context, err error) {
	c.JSON(http.StatusOK, HTTPResponse{
		Code:    http.StatusOK,
		Message: err.Error(),
		Data:    nil,
	})
}
func GinSuccess(c *gin.Context, data any) {
	c.JSON(http.StatusOK, HTTPResponse{
		Code:    http.StatusOK,
		Message: "ok",
		Data:    data,
	})
}

func FromHTTPResponse[T any](response *HTTPResponse) *T {
	bs, _ := json.Marshal(response.Data)
	var res = new(T)
	_ = json.Unmarshal(bs, res)
	return res
}
