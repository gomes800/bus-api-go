package service

import (
	"context"
	"log"
)

func (s *Service) UpdateData() error {
	ctx := context.Background()

	buses, err := s.FetchData()
	if err != nil {
		return err
	}

	log.Printf("[Updater] upstream returned %d buses", len(buses))

	saved := 0
	for _, bus := range buses {
		if err := s.Repo.SaveBus(ctx, bus); err != nil {
			log.Printf("[Updater] SaveBus failed ordem=%s linha=%s: %v", bus.Ordem, bus.Linha, err)
			continue
		}
		saved++
	}

	log.Printf("[Updater] saved %d/%d buses", saved, len(buses))
	return nil
}
