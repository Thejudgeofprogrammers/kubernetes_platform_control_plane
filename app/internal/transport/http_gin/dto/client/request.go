package dto

type ListClientQuery struct {
	Status string `form:"status"`
	Limit  int    `form:"limit" binding:"gte=0,lte=100"`
	Offset int    `form:"offset" binding:"gte=0"`
}

type CreateClientRequest struct {
	Name         string `json:"name" binding:"required"`
	APIServiceID string `json:"api_service_id" binding:"required"`
	Description  string `json:"description,omitempty"`
}

type RestartClientRequest struct {
	Reason string `json:"reason,omitempty"`
}
