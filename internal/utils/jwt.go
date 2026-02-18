package utils

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const (
	// ContextKeyClaims is the key used to store JWT claims in the Gin context.
	ContextKeyClaims = "claims"

	authorizationHeader = "Authorization"
	bearerPrefix        = "Bearer "
)

// ExtractBearerToken returns the token from "Bearer <token>" or an error if missing/invalid.
func ExtractBearerToken(authHeader string) (string, error) {
	if authHeader == "" {
		return "", errors.New("Authorization header is required")
	}
	if !strings.HasPrefix(authHeader, bearerPrefix) {
		return "", errors.New("Authorization header must start with 'Bearer '")
	}
	token := strings.TrimPrefix(authHeader, bearerPrefix)
	if token == "" {
		return "", errors.New("Token is required")
	}
	return token, nil
}

// ParseAndValidate parses a JWT string and validates it with the given secret.
// Returns the claims or an error.
func ParseAndValidate(tokenString string, secret string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid or expired token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}
	return claims, nil
}

// GetClaimsFromContext returns the JWT claims from the Gin context, or (nil, false) if missing/invalid.
func GetClaimsFromContext(c *gin.Context) (jwt.MapClaims, bool) {
	raw, exists := c.Get(ContextKeyClaims)
	if !exists {
		return nil, false
	}
	claims, ok := raw.(jwt.MapClaims)
	return claims, ok
}

// GetUserIDFromContext returns the user_id (JWT "sub" claim) from the request context, or "" if missing.
func GetUserIDFromContext(c *gin.Context) string {
	claims, ok := GetClaimsFromContext(c)
	if !ok {
		return ""
	}
	if sub, ok := claims["sub"].(string); ok {
		return sub
	}
	return ""
}
