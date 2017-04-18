package jsonrpc

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
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

func HandleRPCRequest(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		renderResponse(NewError(err, InvalidRequestErrorCode), w)
		return
	}
	var req Request
	err = json.Unmarshal(body, &req)
	if err != nil {
		renderResponse(NewError(err, ParseErrorCode), w)
		return
	}

	response, errorResponse := HandleRequest(&req)
	if errorResponse != nil {
		renderResponse(NewError(err, InternalErrorCode), w)
		return
	}
	renderResponse(response, w)
	return
}

func renderResponse(res interface{}, w http.ResponseWriter) {
	response, err := json.Marshal(res)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

