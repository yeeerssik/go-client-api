package models

import (
	"go_client_service/core/helpers"
	"time"

	"github.com/google/uuid"
)

type Client struct {
	ID        *int       `gorm:"Column:id" sql:"type:int;not null"`
	ClientID  *[]byte    `gorm:"Column:client_id;type:uuid" sql:"type:binary(16);not null"`
	FirstName *string    `gorm:"Column:first_name" sql:"type:varchar(50);default:null"`
	LastName  *string    `gorm:"Column:last_name" sql:"type:varchar(50);default:null"`
	BirthDate *time.Time `gorm:"Column:birth_date" sql:"type:date;default:null"`
	IIN       *string    `gorm:"Column:iin" sql:"type:varchar(12);not null"`
}

func NewClient() *Client {
	return &Client{}
}

func (*Client) TableName() string {
	return "clients"
}

func CreateClient(
	firstName, lastName *string, birthDate *time.Time, iin *string,
) (
	c *Client,
	err error,
) {
	id, _ := uuid.New().MarshalBinary()
	c = &Client{
		ClientID:  &id,
		FirstName: firstName,
		LastName:  lastName,
		BirthDate: birthDate,
		IIN:       iin,
	}
	if insert := db.Create(c); insert.Error != nil {
		return nil, insert.Error
	}
	return
}

// TODO: Get all clients from table
func GetAllClients() ([]*Client, error) {
	var c []*Client
	if clients := db.Find(&c); clients.Error != nil {
		return nil, clients.Error
	}
	return c, nil
}

func GetClientById(id string) (*Client, error) {
	idInBinary, err := helpers.StringToUUIDByte(id)
	if err != nil {
		return nil, err
	}
	var c Client
	c.ClientID = &idInBinary
	if find := db.First(&c); find.Error != nil {
		return nil, find.Error
	}
	return &c, nil
}

func DeleteClientById(id string) (err error) {
	idInBinary, err := helpers.StringToUUIDByte(id)
	if err != nil {
		return err
	}
	var c Client
	c.ClientID = &idInBinary
	delete := db.Delete(&c, 1)
	if delete.Error != nil {
		return delete.Error
	}
	return nil
}

func UpdateClient(c *Client) (err error) {
	if update := db.Save(c); update.Error != nil {
		return update.Error
	}
	return nil
}
