// middleware.go
package middleware

import (
	"net/http"

	"los-complejos-backend/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware validates the JWT and extracts the user's role, username, and ID
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the token from the Authorization header
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": "Authorization token is required",
			})
			c.Abort()
			return
		}

		// Parse and validate the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Ensure the token uses the correct signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return utils.JWTSecret, nil
		})

		// Handle parsing or validation errors
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": "Invalid or expired token",
			})
			c.Abort()
			return
		}

		// Extract claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": "Invalid token claims",
			})
			c.Abort()
			return
		}

		// Extract and validate required claims
		role, roleOk := claims["role"].(string)
		username, usernameOk := claims["username"].(string)
		id, idOk := claims["_id"].(string)

		if !roleOk || role == "" {
			c.JSON(http.StatusForbidden, gin.H{
				"status":  "error",
				"message": "Role is missing or invalid in the token",
			})
			c.Abort()
			return
		}

		if !usernameOk || username == "" {
			c.JSON(http.StatusForbidden, gin.H{
				"status":  "error",
				"message": "Username is missing or invalid in the token",
			})
			c.Abort()
			return
		}

		if !idOk || id == "" {
			c.JSON(http.StatusForbidden, gin.H{
				"status":  "error",
				"message": "User ID is missing or invalid in the token",
			})
			c.Abort()
			return
		}

		// Store values in the Gin context for downstream handlers
		utils.SetContextValues(c, map[string]interface{}{
			"_id":      id,
			"username": username,
			"role":     role,
		})

		// Proceed to the next handler
		c.Next()
	}
}
