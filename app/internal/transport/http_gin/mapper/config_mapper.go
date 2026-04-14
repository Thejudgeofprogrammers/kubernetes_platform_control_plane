package mapper

import (
	"control_plane/internal/domain"
	configDTO "control_plane/internal/transport/http_gin/dto/config"
	"time"
)

func ToConfigResponse(c *domain.APIClientConfig) configDTO.ConfigResponse {
	return configDTO.ConfigResponse{
		ID:           c.ID,
		ClientID:     c.ClientID,
		Version:      c.Version,
		AuthType:     string(c.AuthType),
		AuthRef:      c.AuthRef,
		TimeoutMs:    c.TimeoutMs,
		RetryCount:   c.RetryCount,
		RetryBackoff: c.RetryBackoff,
		Headers:      c.Headers,
		CreatedAt:    c.CreatedAt.Format(time.RFC3339),
		CreatedBy:    c.CreatedBy,
	}
}
