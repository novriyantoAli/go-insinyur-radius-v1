package domain

import "github.com/dgrijalva/jwt-go"

// JWTCustomClaims ...
type JWTCustomClaims struct {
	Username string `json:"username"`
	ID       uint   `json:"id"`
	Token    string `json:"token"`
	jwt.StandardClaims
}
