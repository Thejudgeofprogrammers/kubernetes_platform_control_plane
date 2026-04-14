package dto

type APIServiceResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	BaseURL   string `json:"base_url"`
	Protocol  string `json:"protocol"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
}
