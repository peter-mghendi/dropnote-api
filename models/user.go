package models

import (
	"github.com/dgrijalva/jwt-go"
	uuid "github.com/satori/go.uuid"
)

// Token is a JWT claims struct
type Token struct {
	UserID uuid.UUID
	jwt.StandardClaims
}

// Account is a struct to rep user account
type Account struct {
	Base
	Name string `json:"email"`
	Pass string `json:"password"`
	Auth string `sql:"-" json:"token"`
}
