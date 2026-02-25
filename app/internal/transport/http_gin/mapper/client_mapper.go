package mapper

import (
	"control_plane/internal/domain"
	dto "control_plane/internal/transport/http_gin/dto/client"
)

func ToClientSummary(c domain.APIClient) dto.ClientSummaryResponse {
	return dto.ClientSummaryResponse{
		ID:           c.ID,
		Name:         c.Name,
		APIServiceID: c.APIServiceID,
		Status:       string(c.GetStatus()),
	}
}

func ToClientSummaryList(list []domain.APIClient) []dto.ClientSummaryResponse {
	result := make([]dto.ClientSummaryResponse, 0, len(list))

	for _, client := range list {
		result = append(result, ToClientSummary(client))
	}

	return result
}
 