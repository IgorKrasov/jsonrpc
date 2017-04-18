package main

import (
	"net/http"
	"github.com/labstack/echo"
	"github.com/IgorKrasov/jsonrpc"
	"encoding/json"
)

type testParams struct {
	Hello string `json:"hello"`
}

func HandleRPCRequest(c echo.Context) error {
	request := new(jsonrpc.Request)
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusInternalServerError, jsonrpc.NewError(err, jsonrpc.InvalidRequestErrorCode))
	}

	result, errorResponse := jsonrpc.HandleRequest(request)
	if errorResponse != nil {
		return c.JSON(http.StatusInternalServerError, errorResponse)
	}

	return  c.JSON(http.StatusOK, result)
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
	e := echo.New()
	e.POST("/rpc", HandleRPCRequest)
	e.Logger.Fatal(e.Start(":3000"))
}
