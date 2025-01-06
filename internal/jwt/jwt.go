package jwt

import (
	"time"
	"webapp/internal/config"
	"webapp/internal/models"

	"github.com/golang-jwt/jwt"
)

func GenerateToken(user *models.User, cfg *config.Config) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(time.Minute * time.Duration(cfg.Jwt.ExpTime)).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(cfg.Jwt.Secret))

	if err != nil {
		return "", err
	}

	return signedToken, nil
}
