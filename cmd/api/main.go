package main

import (
	"log"
	"net/http"

	"github.com/gomes800/bus-api-go/internal/api"
	"github.com/gomes800/bus-api-go/internal/config"
	"github.com/gomes800/bus-api-go/internal/jobs"
	"github.com/gomes800/bus-api-go/internal/service"
)

func main() {
	cfg := config.Load()

	srv := service.New(cfg.Redis)
	handler := api.NewHandler(srv)

	jobs.StartUpdaterJob(srv)

	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	log.Println("API running on :8080")
	http.ListenAndServe(":8080", mux)
}
