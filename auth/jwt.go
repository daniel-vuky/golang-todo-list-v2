package auth

import (
	"fmt"
	jwtGo "github.com/golang-jwt/jwt/v5"
	"os"
	"time"
)

// Create JWT token base on username
func Create(username string) (string, error) {
	claims := jwtGo.MapClaims{}
	claims["authorized"] = true
	claims["user_name"] = username
	claims["exp"] = time.Now().Add(time.Hour).Unix()
	token := jwtGo.NewWithClaims(jwtGo.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("SECRET_JWT")))
}

// ValidateToken pass from API
func ValidateToken(tokenString string) (*jwtGo.Token, error) {
	finalToken, err := jwtGo.Parse(tokenString, func(token *jwtGo.Token) (interface{}, error) {
		// Check the signing method
		if _, ok := token.Method.(*jwtGo.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("SECRET_JWT")), nil
	})

	if err != nil {
		return nil, err
	}

	return finalToken, nil
}

// GetUsernameFromToken get username from token
func GetUsernameFromToken(tokenString string) (string, error) {
	token, tokenErr := ValidateToken(tokenString)

	if tokenErr != nil {
		return "", nil
	}
	if claims, ok := token.Claims.(jwtGo.MapClaims); ok && token.Valid {
		if username, exists := claims["user_name"]; exists {
			return username.(string), nil
		}
	}
	return "", fmt.Errorf("Invalid Token!")
}

// GetUserIDFromToken get username from token
func GetUserIDFromToken(tokenString string) (uint64, error) {
	token, tokenErr := ValidateToken(tokenString)

	if tokenErr != nil {
		return 0, nil
	}
	if claims, ok := token.Claims.(jwtGo.MapClaims); ok && token.Valid {
		if userID, exists := claims["user_id"]; exists {
			return userID.(uint64), nil
		}
	}
	return 0, fmt.Errorf("Invalid Token!")
}
