package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"brainbash_backend/internal/utils"
)

// AuthMiddleware returns a Gin middleware that validates JWT tokens
// from the Authorization header. On success, it stores the parsed
// claims in the context under the key "claims".
func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		tokenString, err := utils.ExtractBearerToken(authHeader)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		claims, err := utils.ParseAndValidate(tokenString, jwtSecret)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		c.Set(utils.ContextKeyClaims, claims)
		c.Next()
	}
}
