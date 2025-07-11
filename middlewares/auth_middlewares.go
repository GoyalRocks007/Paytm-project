package middlewares

import (
	"log"
	"net/http"
	authmodule "paytm-project/internal/modules/auth_module"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the Authorization header
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			c.Abort()
			return
		}

		// Expected format: Bearer <token>
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
			c.Abort()
			return
		}

		token := parts[1]

		claims, err := authmodule.VerifyJWT(token)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		c.Set("claims", claims)

		c.Next() // Proceed to the next middleware or controller
	}
}

func OtpAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the Authorization header
		authHeader := c.GetHeader("two-fac-auth")

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Otp header missing"})
			c.Abort()
			return
		}

		// Expected format: Bearer <token>
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
			c.Abort()
			return
		}

		token := parts[1]

		claims, err := authmodule.VerifyJWT(token)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		c.Set("otp_claims", claims)

		c.Next() // Proceed to the next middleware or controller
	}
}

func CheckAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, ok := authmodule.GetClaims(c)
		if !ok {
			log.Println("No claims present!")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		if claims["role"] != string(authmodule.AdminRoleName) {
			log.Println("user is not admin")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		if email, ok := claims["email"].(string); ok {
			c.Set("email", email)
		} else {
			log.Println("Email claim not found")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		c.Next() // Proceed to the next middleware or controller
	}
}
