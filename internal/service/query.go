package service

import (
	"context"

	"github.com/gomes800/bus-api-go/internal/model"
)

func (s *Service) GetLine(line string) ([]model.BusPosition, error) {
	return s.Repo.GetLine(context.Background(), line)
}

func (s *Service) GetBus(ordem string) (*model.BusPosition, error) {
	return s.Repo.GetBus(context.Background(), ordem)
}
