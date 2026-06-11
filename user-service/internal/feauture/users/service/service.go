package service

import (
	"context"
	core_domain "user-service/internal/core/domain"
	core_logger "user-service/internal/core/logger"
)


type Service struct {
	repo Repo
	log *core_logger.Logger
}

type Repo interface {
	SaveUser(
		ctx context.Context,
		user core_domain.SaveUser,
	) (error)
}

func NewService(
	repo Repo,
) *Service {
	return &Service{
		repo: repo,
	}
}
