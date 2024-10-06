package utils

import (
	"fmt"
	"time"

	"system.buon18.com/m/config"

	"github.com/golang-jwt/jwt/v5"
)

type WebTokenClaims struct {
	Email       string   `json:"email"`
	Role        string   `json:"role"`
	Permissions []string `json:"permissions"`
	jwt.Claims
}

type RefreshTokenClaims struct {
	Email string `json:"email"`
	jwt.Claims
}

func RemoveBearer(token string) (string, error) {
	const BEARER_SCHEMA = "Bearer "
	if token == "" {
		return "", fmt.Errorf("token is required")
	}

	// Remove Bearer schema
	token = token[len(BEARER_SCHEMA):]

	return token, nil
}

func GenerateWebToken(c WebTokenClaims) (string, error) {
	config := config.GetConfigInstance()

	// Create the Claims
	claims := &jwt.MapClaims{
		"email":       c.Email,
		"role":        c.Role,
		"permissions": c.Permissions,
		"exp":         time.Now().Add(time.Second * time.Duration(config.TOKEN_DURATION_SEC)).Unix(),
	}

	// Generate token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token
	tokenString, err := token.SignedString([]byte(config.TOKEN_KEY))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func GenerateRefreshToken(c RefreshTokenClaims) (string, error) {
	config := config.GetConfigInstance()

	// Create the Claims
	claims := &jwt.MapClaims{
		"email": c.Email,
		"exp":   time.Now().Add(time.Duration(config.REFRESH_TOKEN_SEC) * time.Second).Unix(),
	}

	// Generate token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token
	tokenString, err := token.SignedString([]byte(config.REFRESH_TOKEN_KEY))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateWebToken(tokenString string) (WebTokenClaims, error) {
	config := config.GetConfigInstance()

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.TOKEN_KEY), nil
	})
	if err != nil {
		return WebTokenClaims{}, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		permissions, ok := claims["permissions"]
		if !ok {
			return WebTokenClaims{}, fmt.Errorf("unable to get permissions")
		}

		permissionsArr := make([]string, 0)
		for _, permission := range permissions.([]interface{}) {
			permissionsArr = append(permissionsArr, permission.(string))
		}

		return WebTokenClaims{
			Email:       claims["email"].(string),
			Role:        claims["role"].(string),
			Permissions: permissionsArr,
		}, nil
	}

	return WebTokenClaims{}, fmt.Errorf("invalid token")
}

func ValidateRefreshToken(tokenString string) (RefreshTokenClaims, error) {
	config := config.GetConfigInstance()

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.REFRESH_TOKEN_KEY), nil
	})

	if err != nil {
		return RefreshTokenClaims{}, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return RefreshTokenClaims{
			Email: claims["email"].(string),
		}, nil
	}

	return RefreshTokenClaims{}, fmt.Errorf("invalid token")
}
