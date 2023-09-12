package contracts

type BaseRequest struct {
	Method    *string `json:"-"`
	RequestID *string `json:"-"`
}

type BaseResponse struct {
	RequestID *string    `json:"request_id"`
	Method    *string    `json:"method"`
	HTTPCode  *int       `json:"http_code"`
	ErrorData *ErrorData `json:"error_data"`
}

type ErrorData struct {
	Code        uint64 `json:"code"`
	Description string `json:"description"`
}

func Init() error {
	return nil
}
