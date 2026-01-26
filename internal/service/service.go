package service

import (
	"context"

	"github.com/gomes800/bus-api-go/internal/repository"
	"github.com/redis/go-redis/v9"
)

type Service struct {
	Repo *repository.RedisRepo
}

func New(redis *redis.Client) *Service {
	return &Service{
		Repo: repository.NewRedisRepo(redis),
	}
}

func (s *Service) FetchAndUpdate(ctx context.Context) error {
	buses, err := s.FetchData()
	if err != nil {
		return err
	}

	for _, b := range buses {
		if err := s.Repo.SaveBus(ctx, b); err != nil {
			return err
		}
	}

	return nil
}
