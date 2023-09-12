package handlers

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"go_client_service/contracts"
	"go_client_service/core/helpers"
	"go_client_service/core/middlewares"
	"go_client_service/models"

	"github.com/google/uuid"
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	"github.com/scsbatu/go-api/utils/ierror"
	"gorm.io/gorm"
)

type ClientHandler struct{}

func (handler ClientHandler) Any(c echo.Context) error {
	switch c.Get("Method").(string) {
	case http.MethodPost:
		return handler.Post(c)

	case http.MethodGet:
		return handler.Get(c)

	case http.MethodPut:
		return handler.Update(c)
	}
	return RawResponse(c, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
}

func (ClientHandler) Post(c echo.Context) error {
	var response contracts.CreateClientResponse
	var responseCode int
	requestID := c.Get("RequestID").(string)
	method := c.Get("Method").(string)
	response.Method = &method
	response.RequestID = &requestID
	c.Set("path", "client")
	c.Set("method", strings.ToLower(method))
	req := new(contracts.CreateClientRequest)
	if err := helpers.ExtractAndValidate(c, req); err != nil {
		responseCode = http.StatusBadRequest
		response.HTTPCode = &responseCode
		response.ErrorData = err
		e, _ := json.Marshal(err)
		log.Error(c, "[client] ExtractAndValidate error", string(e))
		return RawResponse(c, response, responseCode)
	}
	rq, _ := json.Marshal(req)
	log.Info(c, "[client] Got Req: ", string(rq))
	responseCode, err := createClient(c, req, &response)
	if responseCode == 0 {
		responseCode = http.StatusInternalServerError
	}
	if err != nil {
		e := &contracts.ErrorData{
			Code:        err.GetCode(),
			Description: err.Error(),
		}
		response.ErrorData = e
		log.Error(c, "[client] Got Error:", err)
	}
	response.HTTPCode = &responseCode
	rs, _ := json.Marshal(response)
	log.Info(c, "[client] Response: ", string(rs))
	return RawResponse(c, response, responseCode)
}

func (ClientHandler) Get(c echo.Context) error {
	var response contracts.GetClientResponse
	var responseCode int
	requestID := c.Get("RequestID").(string)
	method := c.Get("Method").(string)
	response.Method = &method
	response.RequestID = &requestID
	c.Set("path", "client")
	c.Set("method", strings.ToLower(method))
	req := new(contracts.GetClientRequest)
	clientID := c.Param("client_id")
	fmt.Println("params", c.ParamNames(), c.ParamValues())
	if clientID != "" {
		req.ClientID = &clientID
	}
	if err := helpers.ExtractAndValidate(c, req); err != nil {
		responseCode = http.StatusBadRequest
		response.HTTPCode = &responseCode
		response.ErrorData = err
		e, _ := json.Marshal(err)
		log.Error(c, "[client] ExtractAndValidate error", string(e))
		return RawResponse(c, response, responseCode)
	}
	rq, _ := json.Marshal(req)
	log.Info(c, "[client] Got Req: ", string(rq))
	responseCode, err := fetchClient(c, req, &response)
	if responseCode == 0 {
		responseCode = http.StatusInternalServerError
	}
	if err != nil {
		e := &contracts.ErrorData{
			Code:        err.GetCode(),
			Description: err.Error(),
		}
		response.ErrorData = e
		log.Error(c, "[client] Got Error:", err)
	}
	response.HTTPCode = &responseCode
	rs, _ := json.Marshal(response)
	log.Info(c, "[client] Response: ", string(rs))
	return RawResponse(c, response, responseCode)
}

func (ClientHandler) Update(c echo.Context) error {
	var response contracts.UpdateClientResponse
	var responseCode int
	requestID := c.Get("RequestID").(string)
	method := c.Get("Method").(string)
	response.Method = &method
	response.RequestID = &requestID
	c.Set("path", "client")
	c.Set("method", strings.ToLower(method))
	req := new(contracts.UpdateClientRequest)
	clientID := c.Param("client_id")
	if clientID != "" {
		req.ClientID = &clientID
	}
	if err := helpers.ExtractAndValidate(c, req); err != nil {
		responseCode = http.StatusBadRequest
		response.HTTPCode = &responseCode
		response.ErrorData = err
		e, _ := json.Marshal(err)
		log.Error(c, "[client] ExtractAndValidate error", string(e))
		return RawResponse(c, response, responseCode)
	}
	rq, _ := json.Marshal(req)
	log.Info(c, "[client] Got Req: ", string(rq))
	responseCode, err := updateClient(c, req, &response)
	if responseCode == 0 {
		responseCode = http.StatusInternalServerError
	}
	if err != nil {
		e := &contracts.ErrorData{
			Code:        err.GetCode(),
			Description: err.Error(),
		}
		response.ErrorData = e
		log.Error(c, "[client] Got Error:", err)
	}
	response.HTTPCode = &responseCode
	rs, _ := json.Marshal(response)
	log.Info(c, "[client] Response: ", string(rs))
	return RawResponse(c, response, responseCode)
}

func createClient(
	c echo.Context,
	req *contracts.CreateClientRequest,
	resp *contracts.CreateClientResponse,
) (
	int,
	ierror.IError,
) {
	cl, err := models.CreateClient(req.FirstName, req.LastName, req.BirthDate, req.IIN)
	if err != nil {
		return http.StatusInternalServerError, middlewares.ErrStatusInternalServerError("Database error", err)
	}
	id := hex.EncodeToString(*cl.ClientID)
	idInStr := uuid.Must(uuid.Parse(id)).String()
	d := contracts.CreateClientData{
		ClientID: &idInStr,
	}
	resp.Data = &d
	return http.StatusOK, nil
}

func fetchClient(
	c echo.Context,
	req *contracts.GetClientRequest,
	resp *contracts.GetClientResponse,
) (
	int,
	ierror.IError,
) {
	cl, err := models.GetClientById(*req.ClientID)
	if err != nil && err != gorm.ErrRecordNotFound {
		return http.StatusInternalServerError, middlewares.ErrStatusInternalServerError("Database error", err)
	} else if err != nil {
		return http.StatusNotFound, middlewares.ErrStatusNotFound("record not found")
	}
	id := hex.EncodeToString(*cl.ClientID)
	idInStr := uuid.Must(uuid.Parse(id)).String()

	d := contracts.GetClientData{
		ClientID:  &idInStr,
		FirstName: cl.FirstName,
		LastName:  cl.LastName,
		BirthDate: cl.BirthDate,
		IIN:       cl.IIN,
	}
	resp.Data = &d
	return http.StatusOK, nil
}

func updateClient(
	c echo.Context,
	req *contracts.UpdateClientRequest,
	resp *contracts.UpdateClientResponse,
) (
	int,
	ierror.IError,
) {
	cl, err := models.GetClientById(*req.ClientID)
	if err != nil && err != gorm.ErrRecordNotFound {
		return http.StatusInternalServerError, middlewares.ErrStatusInternalServerError("Database error", err)
	} else if err != nil {
		return http.StatusNotFound, middlewares.ErrStatusNotFound("record not found")
	}
	if req.FirstName != nil {
		cl.FirstName = req.FirstName
	}
	if req.LastName != nil {
		cl.LastName = req.LastName
	}
	if req.BirthDate != nil {
		cl.BirthDate = req.BirthDate
	}
	if req.IIN != nil {
		cl.IIN = req.IIN
	}
	if err = models.UpdateClient(cl); err != nil {
		return http.StatusInternalServerError, middlewares.ErrStatusInternalServerError("Database error", err)
	}
	success := "true"
	d := contracts.UpdateClientData{
		Success: &success,
	}
	resp.Data = &d
	return http.StatusOK, nil
}
