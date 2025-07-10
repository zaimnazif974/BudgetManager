package utils

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/zaimnazif974/budgeting-BE/pkg/config"
)

type JWTClaims struct {
	UserID uuid.UUID `json:"user_id"`
	Role   string    `json:"role"`
	jwt.RegisteredClaims
}

var jwtKey = []byte(config.GetEnv("JWT_SECRET_KEY", "nil"))

func JwtKey() []byte {
	return jwtKey
}
