package middleware

import (
	"time"

	"github.com/golang-jwt/jwt/v5"

	"auth-service/config"
)

func GenerateAccessToken(userID, name, userableType string) (string, error) {
	claims := jwt.MapClaims{
		"user_id":       userID,
		"name":          name,
		"userable_type": userableType,
		"exp":           time.Now().Add(time.Duration(config.AppConfig.AccessTokenExpMin) * time.Minute).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.AppConfig.AccessTokenSecret))
}

func GenerateRefreshToken(userID, name, userableType string) (string, time.Time, error) {
	exp := time.Now().Add(time.Duration(config.AppConfig.RefreshTokenExpDays) * 24 * time.Hour)
	claims := jwt.MapClaims{
		"user_id":       userID,
		"name":          name,
		"userable_type": userableType,
		"exp":           exp.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(config.AppConfig.RefreshTokenSecret))
	return signed, exp, err
}

func ParseAccessToken(token string) (jwt.MapClaims, error) {
	parsed, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte(config.AppConfig.AccessTokenSecret), nil
	})
	if err != nil || !parsed.Valid {
		return nil, err
	}
	return parsed.Claims.(jwt.MapClaims), nil
}

func ParseRefreshToken(token string) (jwt.MapClaims, error) {
	parsed, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte(config.AppConfig.RefreshTokenSecret), nil
	})
	if err != nil || !parsed.Valid {
		return nil, err
	}
	return parsed.Claims.(jwt.MapClaims), nil
}