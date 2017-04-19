package jsonrpc

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

const (
	Version = "2.0"
)

type endpoint struct {
	Services Services

}

type Request struct {
	Version string          `json:"jsonrpc"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params"`
	Id      int             `json:"id"`
}

type successResponse struct {
	Version string      `json:"jsonrpc"`
	Result  interface{} `json:"result"`
	Id      int         `json:"id"`
}

type Procedure func(json.RawMessage) (result interface{}, err error)

type Services map[string]Procedure


func New() *endpoint {
	return &endpoint{
		Services: map[string]Procedure{},
	}
}

func (e *endpoint) RegisterProcedure(title string, procedure Procedure) {
	e.Services[title] = procedure
}

func (r Request) isValid() bool {
	if r.Version != Version {
		return false
	}

	if r.Method == "" {
		return false
	}

	return true
}

func (e *endpoint) HandleRequest(request *Request) (*successResponse, *errorResponse) {
	if !request.isValid() {
		return nil, e.NewError(nil, InvalidRequestErrorCode)
	}

	procedure, ok := e.Services[request.Method]
	if !ok {
		return nil, e.NewError(nil, MethodNotFoundErrorCode)
	}

	result, err := procedure(request.Params)
	if err != nil {
		return nil, e.NewError(err, InternalErrorCode)
	}

	return &successResponse{
		Version: Version,
		Result:  result,
		Id:      request.Id,
	}, nil
}

func (e *endpoint) HandleRPCRequest(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		e.renderResponse(e.NewError(err, InvalidRequestErrorCode), w)
		return
	}
	var req Request
	err = json.Unmarshal(body, &req)
	if err != nil {
		e.renderResponse(e.NewError(err, ParseErrorCode), w)
		return
	}

	response, errorResponse := e.HandleRequest(&req)
	if errorResponse != nil {
		e.renderResponse(errorResponse, w)
		return
	}
	e.renderResponse(response, w)
	return
}

func (e *endpoint) renderResponse(res interface{}, w http.ResponseWriter) {
	response, err := json.Marshal(res)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

