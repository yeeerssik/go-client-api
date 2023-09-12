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
	Data *UpdateClientData `json:"data"`
}

type UpdateClientData struct {
	Success *string `json:"success"`
}