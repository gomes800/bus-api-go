package jobs

import (
	"log"
	"time"

	"github.com/gomes800/bus-api-go/internal/service"
)

func StartUpdaterJob(s *service.Service) {
	go func() {
		ticker := time.NewTicker(25 * time.Second)
		for range ticker.C {
			log.Println("[Updater] tick")
			if err := s.UpdateData(); err != nil {
				log.Println("[Updater] error:", err)
			}
		}
	}()
}
