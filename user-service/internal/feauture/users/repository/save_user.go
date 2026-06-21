package repository

import (
	"context"
	core_domain "user-service/internal/core/domain"

	"go.uber.org/zap"
)


func (r *Repo) SaveUser(
	ctx context.Context,
	user core_domain.SaveUser,
) error {
	op := "User.Repo.Register"

	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := 
	`
		INSERT INTO userservice.users (
			id, version, username, phone_number
		)
		VALUES ($1, $2, $3, $4)
		RETURNING id, version, username, phone_number; 
	`

	row := r.pool.QueryRow(
		ctx, 
		query,
		user.ID,
		user.Version,
		user.Username,
		user.PhoneNumber,
	)

	var userModel UserModel
	if err := userModel.Scan(row); err != nil {
		r.log.Error("Failed to save user", zap.String("op", op), zap.Error(err))
		return err
	}

	r.log.Debug("User saved successfully:", zap.String("op", op))
	return nil
}