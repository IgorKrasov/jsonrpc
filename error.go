package jsonrpc

const (
	ParseErrorCode          = -32700
	InvalidRequestErrorCode = -32600
	MethodNotFoundErrorCode = -32601
	InvalidParamsErrorCode  = -32602
	InternalErrorCode       = -32603
)

var ErrorMessages = map[int]string{
	ParseErrorCode:          "Parse error",
	InvalidRequestErrorCode: "Invalid Request",
	MethodNotFoundErrorCode: "Method not found",
	InvalidParamsErrorCode:  "Invalid params",
	InternalErrorCode:       "Internal error",
}

type errorResponse struct {
	Version string `json:"jsonrpc"`
	Error   Error  `json:"error"`
	Id      int    `json:"id"`
}

type Error struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func (e *Endpoint) NewError(err error, errorCode int) *errorResponse {
	if err != nil {
		return &errorResponse{
			Version: Version,
			Error: Error{
				Code:    errorCode,
				Message: ErrorMessages[errorCode],
				Data:    err.Error(),
			},
			Id: 1,
		}
	} else {
		return &errorResponse{
			Version: Version,
			Error: Error{
				Code:    errorCode,
				Message: ErrorMessages[errorCode],
			},
			Id: 1,
		}
	}

}
