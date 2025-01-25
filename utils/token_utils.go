// token_utils.go
package utils

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// JWTSecret is the secret key used to sign the tokens.
// Ensure this key is kept secure and not exposed publicly.
var JWTSecret = []byte(os.Getenv("JWT_SECRET"))

// GenerateToken generates a JWT for a user.
// Parameters:
// - id: The user's unique identifier (e.g., database ID).
// - role: The user's role (e.g., "admin" or "user").
// - username: The user's username (e.g., "Xuculup").
// Returns:
// - A signed JWT token as a string.
// - An error if the signing process fails.
func GenerateToken(id, role, username string) (string, error) {
	// Create the claims (payload)
	claims := jwt.MapClaims{
		"_id":      id,       // ID of the user
		"username": username, // Username of the user
		"role":     role,     // Role of the user
	}

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	return token.SignedString(JWTSecret)
}

// SetContextValues sets multiple key-value pairs into the Gin context.
// Parameters:
// - c: The Gin context to which values are added.
// - values: A map of key-value pairs to set in the context.
func SetContextValues(c *gin.Context, values map[string]interface{}) {
	for key, value := range values {
		c.Set(key, value)
	}
}

// Example usage of SetContextValues:
// values := map[string]interface{}{
//     "_id": "12345",
//     "username": "Xuculup",
//     "role": "admin",
// }
// SetContextValues(c, values)
