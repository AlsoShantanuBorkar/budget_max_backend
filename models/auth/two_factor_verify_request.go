package auth

type TwoFactorVerifyRequest struct {
	Code string `json:"code" validate:"required"`
}
