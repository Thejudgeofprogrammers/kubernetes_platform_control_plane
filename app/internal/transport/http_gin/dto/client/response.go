package dto

type ClientSummaryResponse struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	APIServiceID string `json:"api_service_id"`
	Status       string `json:"status"`
}

type ListClientsResponse struct {
	Items []ClientSummaryResponse `json:"items"`
	Total int                     `json:"total"`
}

type ClientResponse struct {
	ID             string  `json:"id"`
	Name           string  `json:"name"`
	Slug         string `json:"slug"`
	URL          string `json:"url"`
	APIServiceID   string  `json:"api_service_id"`
	Status         string  `json:"status"`
	ActiveConfigID *string `json:"active_config_id,omitempty"`
	CreatedAt      string  `json:"created_at"`
}

type ClientActionResponse struct {
	ClientID string `json:"client_id"`
	Action   string `json:"action"`
	Status   string `json:"status"`
}

type DeleteClientResponse struct {
	ClientID string `json:"client_id"`
	Status   string `json:"status"`
}

type ConfigResponse struct {
	ID           string            `json:"id"`
	ClientID     string            `json:"client_id"`
	Version      string            `json:"version"`
	AuthType     string            `json:"auth_type"`
	AuthRef      string            `json:"auth_ref,omitempty"`
	TimeoutMs    int               `json:"timeout_ms"`
	RetryCount   int               `json:"retry_count"`
	RetryBackoff int               `json:"retry_backoff"`
	Headers      map[string]string `json:"headers"`
	CreatedAt    string            `json:"created_at"`
	CreatedBy    string            `json:"created_by"`
}

