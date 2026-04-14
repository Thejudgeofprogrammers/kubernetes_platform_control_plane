package health

import (
	"net/http"
	"time"

	"api-client/internal/config"
)

func Handler(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		client := http.Client{
			Timeout: time.Duration(cfg.TimeoutMs) * time.Millisecond,
		}

		resp, err := client.Get(cfg.BaseURL)
		if err != nil {
			http.Error(w, "unhealthy", 500)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode >= 200 && resp.StatusCode < 500 {
			w.WriteHeader(200)
			w.Write([]byte("ok"))
			return
		}

		http.Error(w, "unhealthy", 500)
	}
}
