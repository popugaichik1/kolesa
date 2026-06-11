package transport_kafka

import (
	"github.com/google/uuid"
)



type UserRegisterEvent struct {
	ID 				uuid.UUID  	`json:"id" validate:"required"`
	Username 		string     	`json:"username" validate:"required,min=1,max=125"`
	PhoneNumber 	string	   	`json:"phone_number" validate:"required,min=10,max-15,startswith=+"`
}

