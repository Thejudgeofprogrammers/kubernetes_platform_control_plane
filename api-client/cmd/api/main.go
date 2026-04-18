package main

import (
	"api-client/internal/config"
	"api-client/internal/health"
	"api-client/internal/proxy"
	"log"
	"net/http"
)

func main() {
	cfg := config.LoadEnv()

	p := proxy.New(cfg)

	http.Handle("/", p)
	
	http.Handle("/health", health.Handler(cfg))

	http.HandleFunc("/metrics", p.MetricsHandler)

	log.Println("api-client started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
