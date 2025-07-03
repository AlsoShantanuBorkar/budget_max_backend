package auth

type RefreshTokensRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required,uuid4"`
}
