package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

const baseUrl = "https://dados.mobilidade.rio/gps/sppo"

type BusData struct {
	Data interface{} `json:"data"`
}

func fetchData() ([]byte, error) {
	now := time.Now()
	start := now.Add(-30 * time.Second)

	layout := "2006-01-02+15:04:05"
	dataInicial := start.Format(layout)
	dataFinal := now.Format(layout)

	url := fmt.Sprintf("%s?dataInicial=%s&dataFinal=%s", baseUrl, dataInicial, dataFinal)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("Query error: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Read error: %w", err)
	}

	return body, nil
}

func busHandler(w http.ResponseWriter, r *http.Request) {
	data, err := fetchData()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func main() {
	http.HandleFunc("/buses", busHandler)
	fmt.Println("Server running on 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
