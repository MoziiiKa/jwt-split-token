package middlewares

import (
	"jwt-split-token/auth"
	"os"

	"github.com/gin-gonic/gin"
)

// it is called when a user send request to access time.ir
func Auth() gin.HandlerFunc {
	return func(context *gin.Context) {

		// extract token from request header
		tokenString := context.GetHeader("Authorization")
		if tokenString == "" {
			context.JSON(401, gin.H{"error": "request does not contain an access token"})
			context.Abort()
			return
		}

		jwtKeyStr := os.Getenv("JWT_KEY")
		jwtKey := []byte(jwtKeyStr)

		// validating the token
		tokenString, err := auth.ValidateToken(tokenString, jwtKey)
		if err != nil {
			context.JSON(403, gin.H{"error": err.Error()})
			context.Abort()
			return
		}

		// set token
		context.Set("token", tokenString)
		context.Next()
	}
}
