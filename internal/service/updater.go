package service

import (
	"context"
)

func (s *Service) UpdateData() error {
	ctx := context.Background()

	buses, err := s.FetchData()
	if err != nil {
		return err
	}

	for _, bus := range buses {
		s.Repo.SaveBus(ctx, bus)
	}

	return nil
}
