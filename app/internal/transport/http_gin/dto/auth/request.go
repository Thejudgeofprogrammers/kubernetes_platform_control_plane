package dto

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email" example:"artemkotov78@mail.ru"`
	FullName string `json:"full_name" binding:"required"`
}

type RequestCodeRequest struct {
	Email string `json:"email" binding:"required,email" example:"artemkotov78@mail.ru"`
}

type VerifyCodeRequest struct {
	Email string `json:"email" binding:"required,email" example:"artemkotov78@mail.ru"`
	Code  string `json:"code" binding:"required"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type UpdateRoleRequest struct {
	Role string `json:"role" binding:"required"`
}
