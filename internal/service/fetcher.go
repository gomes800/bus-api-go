package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/gomes800/bus-api-go/internal/model"
)

const baseURL = "https://dados.mobilidade.rio/gps/sppo"

func (s *Service) FetchData() ([]model.BusPosition, error) {
	now := time.Now()
	start := now.Add(-30 * time.Second)

	u, _ := url.Parse(baseURL)
	q := u.Query()

	q.Set("dataInicial", start.Format(time.RFC3339))
	q.Set("dataFinal", now.Format(time.RFC3339))

	u.RawQuery = q.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("upstream status=%d body=%s", resp.StatusCode, string(body))
	}

	var buses []model.BusPosition
	if err := json.NewDecoder(resp.Body).Decode(&buses); err != nil {
		return nil, err
	}

	return buses, nil
}
