package helpers

import (
	"fmt"
	"go_client_service/contracts"
	"go_client_service/core/middlewares"

	"github.com/go-playground/locales/en_US"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator"
	"github.com/labstack/echo"
)

var en ut.Translator

type CustomValidator struct {
	Validator *validator.Validate
}

func (cv *CustomValidator) Init() {
	english := en_US.New()
	uni := ut.New(english, english)

	en, _ = uni.GetTranslator("en")
	cv.Validator.RegisterTranslation("required", en, func(ut ut.Translator) error {
		return ut.Add("required", "{0} must have a value!", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("required", fe.Field())
		return t
	})
}

func (cv *CustomValidator) Validate(i interface{}) error {
	// cv.Init()
	return cv.Validator.Struct(i)
}

func ExtractAndValidate(c echo.Context, req interface{}) *contracts.ErrorData {
	// Extracting the values
	if err := c.Bind(req); err != nil {
		httpError := middlewares.ErrTypeMismatch()
		e := &contracts.ErrorData{
			Code:        httpError.GetCode(),
			Description: httpError.Error(),
		}
		return e
	}
	// Validate the struct
	if err := c.Validate(req); err != nil {
		errs := err.(validator.ValidationErrors)
		var msg string
		for _, e := range errs {
			msg = fmt.Sprintf("%s %s", msg, e.Translate(en))
		}
		httpError := middlewares.ErrParametersMissing(msg)
		e := &contracts.ErrorData{
			Code:        httpError.GetCode(),
			Description: httpError.Error(),
		}
		return e
	}
	return nil
}
