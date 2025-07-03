package auth

type TwoFactorLoginRequest struct {
	Email string `json:"email" validate:"required,email"`
	Code  string `json:"code" validate:"required"`
	Token string `json:"token" validate:"required"`
}
