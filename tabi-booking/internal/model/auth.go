package model

import "github.com/labstack/echo/v4"

// AuthToken holds authentication token details with refresh token
// swagger:model
type AuthToken struct {
	AccessToken  string `json:"access_token,omitempty"`
	TokenType    string `json:"token_type,omitempty"`
	ExpiresIn    int    `json:"expires_in,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

// RefreshToken  represents data stored in JWT token for refresh token
type RefreshToken struct {
	ExpiredAt int64 `json:"expired_at"`
}

// AuthoPartner represents data stored in JWT token for partner
type AuthoPartner struct {
	ID       int
	Username string
	Email    string
	Role     string
}

type AuthoUser struct {
	ID       int
	Username string
	Email    string
	Phone    string
	Role     string
}

// AuthFile represents auth interface
type Autho interface {
	Partner(echo.Context) *AuthoPartner
	User(echo.Context) *AuthoUser
}

type PartnerTokenClaims struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	Email    string `json:"email"`
}

type UserTokenClaims struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Role     string `json:"role"`
}
