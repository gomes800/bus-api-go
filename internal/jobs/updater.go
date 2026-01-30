package jobs

import (
	"log"
	"time"

	"github.com/gomes800/bus-api-go/internal/service"
)

func StartUpdaterJob(s *service.Service) {
	go func() {
		for {
			log.Println("[Updater] tick")

			if err := s.UpdateData(); err != nil {
				log.Println("[Updater] error:", err)
			}

			time.Sleep(30 * time.Second)
		}
	}()
}
