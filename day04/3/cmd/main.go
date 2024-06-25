package main

import (
	"github.com/gin-gonic/gin"
	"github.com/lahnasti/GO_praktikum/internal/server"
)

func main() {
	r := gin.Default()
	r.POST("/login", server.LoginHandler)
	protected := r.Group("/", server.AuthMiddleware())
	protected.GET("/profile", server.ProfileHandler)

	r.Run(":8080")
}

//R.GROUP - защищены middleware для аутентификации. Маршруты в этой группе требуют
//наличия валидного JWT токена для доступа.

