package dto

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
