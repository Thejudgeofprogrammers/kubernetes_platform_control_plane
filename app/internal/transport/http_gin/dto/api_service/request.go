package dto

type CreateAPIServiceRequest struct {
	Name     string `json:"name" binding:"required"`
	BaseURL  string `json:"base_url" binding:"required"`
	Protocol string `json:"protocol" binding:"required"`
}

type UpdateAPIServiceRequest struct {
	Name     string `json:"name"`
	BaseURL  string `json:"base_url"`
	Protocol string `json:"protocol"`
}
