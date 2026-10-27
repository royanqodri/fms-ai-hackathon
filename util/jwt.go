package util

import (
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/royanqodri/fms-ai-hackathon/config"
)

func GenerateToken(idUser int64, username string, email string, module string, customerNo string, customerId int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"customer_id": customerId,
		"customer_no": customerNo,
		"user_id":     idUser,
		"username":    username,
		"email":       email,
		"module":      module,
		"source":      "web_mobile",
		"exp":         GetTimeNowByLoc().Add(time.Duration(config.Get().JWT.Expire) * time.Second).Unix(),
	})

	return token.SignedString([]byte(config.Get().JWT.SecretKey))
}

func GenerateTokenOnboard(module string, customerNo string, customerId int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"customer_id": customerId,
		"customer_no": customerNo,
		"module":      module,
		"source":      "onboard",
		"exp":         GetTimeNowByLoc().Add(time.Duration(config.Get().JWT.Expire) * time.Second).Unix(),
	})

	return token.SignedString([]byte(config.Get().JWT.SecretKey))
}

// getTokenFromHeader extracts the token from the Authorization header
func GetTokenFromHeader(ctx *gin.Context) string {
	authHeader := ctx.GetHeader("Authorization")
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return ""
	}
	return parts[1]
}

// parseToken validates the JWT token and extracts claims
func ParseToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Get().JWT.SecretKey), nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}

// isTokenExpired checks if the token's expiration time is in the past
func IsTokenExpired(claims jwt.MapClaims) bool {
	expiredTime, ok := claims["exp"].(float64)
	if !ok {
		return true
	}
	return time.Now().After(time.Unix(int64(expiredTime), 0))
}

// isAuthorized checks if the user's role is authorized for the current route
func IsAuthorized(ctx *gin.Context, username string) bool {
	// urlPath := ctx.FullPath()
	// method := ctx.Request.Method
	// rolesNeeded := enum.ApiAccessConfigs[urlPath][method]

	// // Default to authorized if no specific roles are required
	// if len(rolesNeeded) == 0 {
	// 	return true
	// }

	// // Check if the user's role matches any of the allowed roles
	// for _, roleNeeded := range rolesNeeded {
	// 	if constants.Role(role) == roleNeeded {
	// 		return true
	// 	}
	// }

	fmt.Println("username", username)
	return false
}
