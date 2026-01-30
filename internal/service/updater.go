package service

import (
	"context"
	"log"
)

var updateCycle int

func (s *Service) UpdateData() error {
	ctx := context.Background()

	buses, err := s.FetchData()
	if err != nil {
		return err
	}

	log.Printf("[Updater] upstream returned %d buses", len(buses))

	touchedLines := make(map[string]struct{})
	saved := 0

	for _, bus := range buses {
		touchedLines[bus.Linha] = struct{}{}

		if err := s.Repo.SaveBus(ctx, bus); err != nil {
			log.Printf("[Updater] SaveBus failed ordem=%s linha=%s: %v", bus.Ordem, bus.Linha, err)
			continue
		}
		saved++
	}

	log.Printf("[Updater] saved %d/%d buses", saved, len(buses))

	updateCycle++
	if updateCycle%5 == 0 {
		if err := s.Repo.CleanupGeo(ctx, 3000); err != nil {
			log.Printf("[Cleanup] CleanupGeo error: %v", err)
		}

		for line := range touchedLines {
			if err := s.Repo.CleanupLine(ctx, line); err != nil {
				log.Printf("[Cleanup] CleanupLine line=%s error: %v", line, err)
			}
		}

		log.Printf("[Cleanup] done (lines=%d)", len(touchedLines))
	}

	return nil
}
