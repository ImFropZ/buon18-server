package utils

import (
	"fmt"
	"server/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	Email string `json:"email"`
	Role  string `json:"role"`
	jwt.Claims
}

func GenerateWebToken(c Claims) (string, error) {
	config := config.GetAuthConfigInstance()

	// Create the Claims
	claims := &jwt.MapClaims{
		"Email": c.Email,
		"Role":  c.Role,
		"Exp":   time.Now().Add(time.Second * time.Duration(config.TOKEN_DURATION_SEC)).Unix(),
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

func ValidateWebToken(token string) (Claims, error) {
	config := config.GetAuthConfigInstance()

	tokenBytes, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}

		return []byte(config.TOKEN_KEY), nil
	})
	if err != nil {
		return Claims{}, err
	}

	if claims, ok := tokenBytes.Claims.(jwt.MapClaims); !ok && !tokenBytes.Valid {
		return Claims{}, fmt.Errorf("invalid token")
	} else {
		// -- Show human readable time
		expTime := time.Unix(int64(claims["Exp"].(float64)), 0)
		currentTime := time.Now()
		fmt.Println(expTime, currentTime)

		if expTime.Sub(currentTime) < 0 {
			return Claims{}, fmt.Errorf("token expired")
		}

		return Claims{
			Email: claims["Email"].(string),
			Role:  claims["Role"].(string),
		}, nil
	}
}
