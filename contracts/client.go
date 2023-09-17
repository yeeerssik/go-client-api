package contracts

import "time"

type CreateClientRequest struct {
	BaseRequest
	FirstName *string    `json:"first_name" validate:"required"`
	LastName  *string    `json:"last_name"`
	BirthDate *time.Time `json:"birth_date"`
	IIN       *string    `json:"iin"`
}

type CreateClientResponse struct {
	BaseResponse
	Data *CreateClientData `json:"data"`
}

type CreateClientData struct {
	ClientID *string `json:"client_id"`
}

type GetClientRequest struct {
	BaseRequest
	ClientID *string `path:"client_id" validate:"required"`
}

type GetClientResponse struct {
	BaseResponse
	Data *GetClientData `json:"data"`
}

type GetClientData struct {
	ClientID  *string    `json:"client_id" validate:"required"`
	FirstName *string    `json:"first_name" validate:"required"`
	LastName  *string    `json:"last_name"`
	BirthDate *time.Time `json:"birth_date"`
	IIN       *string    `json:"iin"`
}

type GetAllClientRequest struct {
	BaseRequest
}

type GetAllClientResponse struct {
	BaseResponse
	Data []*GetClientData `json:"data"`
}

type DeleteClientRequest struct {
	BaseRequest
	ClientID *string `path:"client_id" validate:"required"`
}

type DeleteClientResponse struct {
	BaseResponse
	Data *ManipulateClientData `json:"data"`
}

type UpdateClientRequest struct {
	BaseRequest
	ClientID  *string    `json:"client_id" validate:"required"`
	FirstName *string    `json:"first_name" validate:"required"`
	LastName  *string    `json:"last_name"`
	BirthDate *time.Time `json:"birth_date"`
	IIN       *string    `json:"iin"`
}

type UpdateClientResponse struct {
	BaseResponse
	Data *ManipulateClientData `json:"data"`
}

type ManipulateClientData struct {
	Success *string `json:"success"`
}
