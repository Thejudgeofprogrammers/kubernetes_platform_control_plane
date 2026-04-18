package impl

import (
	"context"
	"control_plane/internal/domain"
	"control_plane/internal/service/metric"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type metricsResponse struct {
	Requests int64 `json:"requests"`
	Errors   int64 `json:"errors"`
	Latency  int64 `json:"latency"`
}

type metricsService struct {
	client *http.Client
}

func NewMetricsService() metric.MetricsService {
	return &metricsService{
		client: &http.Client{
			Timeout: 3 * time.Second,
		},
	}
}

func (s *metricsService) Collect(
	ctx context.Context,
	baseURL string,
	clientID string,
) ([]domain.Metric, error) {

	url := fmt.Sprintf("%s/api/clients/%s/metrics", baseURL, clientID)

	resp, err := s.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	return ParseMetrics(body, clientID)
}

func ParseMetrics(body []byte, clientID string) ([]domain.Metric, error) {
	var m metricsResponse

	err := json.Unmarshal(body, &m)
	if err != nil {
		return nil, err
	}

	return []domain.Metric{
		{
			ClientID:  clientID,
			Requests:  m.Requests,
			Errors:    m.Errors,
			Latency:   m.Latency,
			CreatedAt: time.Now(),
		},
	}, nil
}
