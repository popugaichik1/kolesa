package repository

import (
	core_domain "user-service/internal/core/domain"
	core_postgres_pool "user-service/internal/core/repository/postgres/pool"

	"github.com/google/uuid"
)

type UserModel struct {
	ID 				uuid.UUID
	Version         int
	Username 		string
	PhoneNumber 	string
}


func (m *UserModel) Scan(
	row core_postgres_pool.Row,
) error {
	return row.Scan(
		&m.ID,
		&m.Version,
		&m.Username,
		&m.PhoneNumber,
	)
}


func modelToDomain(model UserModel) core_domain.User {
	return core_domain.User{
		ID: 			model.ID,
		Version: 		model.Version,
		Username: 	 	model.Username,
		PhoneNumber:  	model.PhoneNumber,
	}
}