package auth

import "github.com/golang-jwt/jwt/v4"

type UserType string

const (
	Admin          UserType = "admin"
	IndividualUser UserType = "individual"
	CorporateUser  UserType = "corporate"
)

type Claims struct {
	Username string   `json:"username"`
	UserType UserType `json:"user_type"`
	UserID   uint
	jwt.RegisteredClaims
}

func (c *Claims) IsAdmin() bool {
	return c.UserType == Admin
}

func (c *Claims) IsNotAdmin() bool {
	return !c.IsAdmin()
}

func (c *Claims) IsIndividualUser() bool {
	return c.UserType == IndividualUser
}

func (c *Claims) IsCorporatedUser() bool {
	return c.UserType == CorporateUser
}
