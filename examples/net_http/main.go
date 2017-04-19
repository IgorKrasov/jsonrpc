package main

import (
	"net/http"
	"encoding/json"
	"github.com/IgorKrasov/jsonrpc"
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
	e := jsonrpc.New()
	e.RegisterProcedure("test.test", testProcedure)
	http.HandleFunc("/rpc", e.HandleRPCRequest)

	http.ListenAndServe(":1323", nil)
}
