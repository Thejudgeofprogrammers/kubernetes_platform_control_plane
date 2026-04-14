package mapper

import (
	"control_plane/internal/domain"
		apiserviceDTO "control_plane/internal/transport/http_gin/dto/api_service"
	"time"
)

func ToAPIServiceResponse(s *domain.APIService) apiserviceDTO.APIServiceResponse {
	return apiserviceDTO.APIServiceResponse{
		ID:        s.ID,
		Name:      s.Name,
		BaseURL:   s.BaseURL,
		Protocol:  s.Protocol,
		Status:    s.Status,
		CreatedAt: s.CreatedAt.Format(time.RFC3339),
	}
}
