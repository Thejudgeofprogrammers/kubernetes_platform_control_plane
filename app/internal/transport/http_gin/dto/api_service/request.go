package dto

type CreateAPIServiceRequest struct {
	Name     string `json:"name" binding:"required"`
	BaseURL  string `json:"base_url" binding:"required"`
	Protocol string `json:"protocol" binding:"required"`
}
