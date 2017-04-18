package main

import (
	"log"
	"net/http"
	"github.com/IgorKrasov/jsonrpc"
	"io/ioutil"
	"encoding/json"
)

type testParams struct {
	Hello string `json:"hello"`
}

func HandleRPCRequest(w http.ResponseWriter, r *http.Request) {
	log.Println("handle request")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		renderResponse(jsonrpc.NewError(err, jsonrpc.InvalidRequestErrorCode), w)
		return
	}
	var req jsonrpc.Request
	err = json.Unmarshal(body, &req)
	if err != nil {
		renderResponse(jsonrpc.NewError(err, jsonrpc.ParseErrorCode), w)
		return
	}

	response, errorResponse := jsonrpc.HandleRequest(&req)
	if errorResponse != nil {
		renderResponse(jsonrpc.NewError(err, jsonrpc.InternalErrorCode), w)
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

func testProcedure(p json.RawMessage) (interface{}, error) {
	var params testParams
	err := json.Unmarshal(p, &params)
	if err != nil {
		return nil, err
	}

	return "hello " + params.Hello, nil
}

func main() {
	jsonrpc.New()
	jsonrpc.RegisterProcedure("test.test", testProcedure)
	http.HandleFunc("/rpc", HandleRPCRequest)

	http.ListenAndServe(":1323", nil)
}
