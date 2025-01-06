package jwt

import (
	"errors"
	"strings"
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

func ValidateToken(tokenString string, cfg *config.Config) (bool, error) {
	if strings.HasPrefix(strings.ToLower(tokenString), "bearer ") {
		tokenString = strings.TrimSpace(strings.TrimPrefix(tokenString, "Bearer "))
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(cfg.Jwt.Secret), nil
	})

	if err != nil {
		return false, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if exp, ok := claims["exp"].(float64); ok {
			if time.Unix(int64(exp), 0).Before(time.Now()) {
				return false, errors.New("token has expired")
			}
		} else {
			return false, errors.New("invalid expiration time")
		}

		return true, nil
	}

	return false, errors.New("invalid token")
}
