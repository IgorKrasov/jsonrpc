package jsonrpc

import (
	"encoding/json"
)

const (
	Version = "2.0"
)

type Request struct {
	Version string          `json:"jsonrpc"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params"`
	Id      int             `json:"id"`
}

type SuccessResponse struct {
	Version string      `json:"jsonrpc"`
	Result  interface{} `json:"result"`
	Id      int         `json:"id"`
}



type Procedure func(json.RawMessage) (result interface{}, err error)

type Services map[string]Procedure

var services Services

func New() *Request {
	services = Services{}

	return &Request{}
}

func RegisterProcedure(title string, procedure Procedure) {
	services[title] = procedure
}

func (r Request) IsValid() bool {
	if r.Version != Version {
		return false
	}

	if r.Method == "" {
		return false
	}

	return true
}

func HandleRequest(request *Request) (*SuccessResponse, *ErrorResponse) {
	// Валидируем запрос
	if !request.IsValid() {
		return nil, NewError(nil, InvalidRequestErrorCode)
	}

	// Ищем процедуру
	procedure, ok := services[request.Method]
	if !ok {
		return nil, NewError(nil, MethodNotFoundErrorCode)
	}

	// Запускаем процедуру с параметрами
	result, err := procedure(request.Params)
	if err != nil {
		return nil, NewError(err, InternalErrorCode)
	}

	return &SuccessResponse{
		Version: Version,
		Result:  result,
		Id:      request.Id,
	}, nil
}

