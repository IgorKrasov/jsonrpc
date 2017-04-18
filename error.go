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

type ErrorResponse struct {
	Version string `json:"jsonrpc"`
	Error   Error  `json:"error"`
	Id      int    `json:"id"`
}

type Error struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func NewError(err error, errorCode int) *ErrorResponse {
	if err != nil {
		return &ErrorResponse{
			Version: Version,
			Error: Error{
				Code:    errorCode,
				Message: ErrorMessages[errorCode],
				Data:    err.Error(),
			},
			Id: 1,
		}
	} else {
		return &ErrorResponse{
			Version: Version,
			Error: Error{
				Code:    errorCode,
				Message: ErrorMessages[errorCode],
			},
			Id: 1,
		}
	}

}
