package handlers

import "github.com/labstack/echo"

func RawResponse(c echo.Context, response interface{}, httpCode int) error {
	var responseFunc func(int, interface{}) error
	switch c.Request().Header.Get("accept") {
	case "application/json", "text/json", "json":
		responseFunc = c.JSON
	case "application/xml", "text/xml", "xml":
		responseFunc = c.XML
	default:
		responseFunc = c.JSON
	}
	return responseFunc(httpCode, response)
}

func Init() error {
	return nil
}
