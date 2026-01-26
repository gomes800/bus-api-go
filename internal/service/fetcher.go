package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gomes800/bus-api-go/internal/model"
)

const baseURL = "https://dados.mobilidade.rio/gps/sppo"

func (s *Service) FetchData() ([]model.BusPosition, error) {
	now := time.Now()
	start := now.Add(-30 * time.Second)

	layout := "2006-01-02+15:04:05"
	url := fmt.Sprintf("%s?dataInicial=%s&dataFinal=%s",
		baseURL, start.Format(layout), now.Format(layout))

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var buses []model.BusPosition
	if err := json.NewDecoder(resp.Body).Decode(&buses); err != nil {
		return nil, err
	}

	return buses, nil
}
