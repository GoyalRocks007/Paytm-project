package authmodule

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

var (
	jwtSecret string
	GetSecret = func() string {
		if jwtSecret != "" {
			return jwtSecret
		}
		jwtSecret = os.Getenv("JWT_SECRET")
		return jwtSecret
	}
)

func SignupRequestDtoToUserMapper(signupRequestDto *SignupRequestDto) *User {
	return &User{
		Name:     signupRequestDto.Name,
		Email:    signupRequestDto.Email,
		Password: signupRequestDto.Password,
		Contact:  signupRequestDto.Contact,
		Username: signupRequestDto.Username,
	}
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost) // DefaultCost = 10
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func CheckPassword(hashedPassword, plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	return err == nil
}

func GenerateJwt(claimsMap map[string]interface{}) (string, error) {
	claims := jwt.MapClaims{}

	// Copy custom claims
	for k, v := range claimsMap {
		claims[k] = v
	}

	now := time.Now()

	// Only set "exp" if it's not already provided
	if _, exists := claims["exp"]; !exists {
		claims["exp"] = now.Add(time.Minute * 30).Unix()
	}

	// Set common registered claims if not already present
	if _, exists := claims["iss"]; !exists {
		claims["iss"] = "paytm-project"
	}

	// Create and sign token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(GetSecret()))
}

func VerifyJWT(tokenString string) (jwt.MapClaims, error) {
	// Parse and validate token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Make sure the signing method is HMAC and specifically HS256
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(GetSecret()), nil
	})

	// Error while parsing
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	// Check if token is valid
	if ok && token.Valid {
		return claims, nil
	}

	// Check expiry
	if exp, ok := claims["exp"].(float64); ok {
		log.Println("hello time 2", exp, time.Now().Unix())
		if int64(exp) < time.Now().Unix() {
			return nil, errors.New("token expired")
		}
	}

	return nil, errors.New("invalid token")
}

func GetClaims(c *gin.Context) (jwt.MapClaims, bool) {
	val, exists := c.Get("claims")
	if !exists {
		return nil, false
	}
	claims, ok := val.(jwt.MapClaims)
	return claims, ok
}
