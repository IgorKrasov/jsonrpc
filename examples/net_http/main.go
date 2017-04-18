package main

import (
	"net/http"
	"encoding/json"
)

type testParams struct {
	Hello string `json:"hello"`
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
	http.HandleFunc("/rpc", jsonrpc.HandleRPCRequest)

	http.ListenAndServe(":1323", nil)
}
