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

type ClientConfigRequest struct {
	Version      string            `json:"version" binding:"required"`
	AuthType     string            `json:"auth_type" binding:"required,oneof=none api_key bearer"`
	AuthRef      string            `json:"auth_ref,omitempty"`
	TimeoutMs    int               `json:"timeout_ms" binding:"gte=0,lte=60000"`
	RetryCount   int               `json:"retry_count" binding:"gte=0,lte=10"`
	RetryBackoff int               `json:"retry_backoff" binding:"gte=0,lte=60000"`
	Headers      map[string]string `json:"headers"`
}
