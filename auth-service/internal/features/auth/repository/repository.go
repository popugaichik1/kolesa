package auth_postgres_repository

import (
	core_logger "github.com/zosinkin/social_network/internal/core/logger"
	core_postgres_pool "github.com/zosinkin/social_network/internal/core/repository/postgres/pool"
)


type UserRepo struct {
	pool core_postgres_pool.Pool
	log  *core_logger.Logger
}

func NewUsersRepo(
	pool core_postgres_pool.Pool,
) *UserRepo {
	return &UserRepo{
		pool: pool,
	}
}


type RefreshTokenRepo struct {
	pool 	core_postgres_pool.Pool
	log  	*core_logger.Logger
}


func NewRefreshTokenRepo(
	pool core_postgres_pool.Pool,
) *RefreshTokenRepo {
	return &RefreshTokenRepo{
		pool: pool,
	}
}