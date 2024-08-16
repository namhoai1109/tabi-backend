package model

import "github.com/labstack/echo/v4"

// AuthoPartner represents data stored in JWT token for partner
type AuthoPartner struct {
	ID       int
	Username string
	Email    string
	Role     string
}

// Autho represents auth interface
type Autho interface {
	Partner(echo.Context) *AuthoPartner
}

const (
	AuthoRoleRepresentative = "REP"
	AuthoRoleBranchManager  = "BMA"
	AuthoRoleClient         = "CLI"
	AuthoRoleHost           = "HST"
)

var PartnerRoles = []string{
	AuthoRoleRepresentative,
	AuthoRoleBranchManager,
	AuthoRoleHost,
}
