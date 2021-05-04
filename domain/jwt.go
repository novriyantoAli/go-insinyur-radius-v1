package domain

import "github.com/dgrijalva/jwt-go"

// JWTCustomClaims ...
type JWTCustomClaims struct {
	Username string  `json:"username"`
	ID       int64   `json:"id"`
	Level    string  `json:"level"`
	Token    *string `json:"token"`
	jwt.StandardClaims
}
