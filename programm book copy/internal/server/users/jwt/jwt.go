package jwt

import (
	"time"
	"net/http"
	"strings"
	"fmt"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/lahnasti/GO_praktikum/internal/models"
)

var jwtSecret = []byte("secure_jwt")

//Функция для генерации JWT токена для пользователя:
func GenerateJWT (user models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims {
		"id": user.ID,
		"name": user.Name,
		"email": user.Email,
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenString, err := token.SignedString(jwtSecret)
		if err != nil {
			return "", err
	}

	return tokenString, nil
}

//Создание JWTAuthMiddleware для проверки токена
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context){
	authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			ctx.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

//Обработчик для входа пользователя, который проверяет учетные данные и генерирует JWT токен:
func LoginHandler(ctx *gin.Context) {
	var user models.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Проверка учетных данных пользователя (например, из базы данных)
	if user.Email == "test@example.com" && user.Password == "password" {
		tokenString, err := GenerateJWT(user)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"token": tokenString})
	} else {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
	}
}