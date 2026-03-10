package jwt

type JWTService interface {
	GenerateAccessToken(userID, role string) (string, error)
	Parse(tokenStr string) (*Claims, error)
}
