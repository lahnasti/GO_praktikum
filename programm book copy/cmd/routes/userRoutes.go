package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/lahnasti/GO_praktikum/internal/server/users"
	"github.com/lahnasti/GO_praktikum/internal/server/users/jwt"
)

func UserRoutes(r *gin.Engine, server *users.Server) {
	userGroup := r.Group("/users")
	{
		userGroup.GET("/", server.GetUsersHandler)
		userGroup.POST("/add", server.RegisterUserHandler)
		userGroup.POST("/adds", server.RegisterMultipleUsersHadler)
		userGroup.POST("/login", server.LoginHandler)

		userGroup.GET("/:id", jwt.JWTAuthMiddleware(), server.GetUserByIDHandler)
		userGroup.PUT("/:id", jwt.JWTAuthMiddleware(), server.UpdateUserHandler)
		userGroup.DELETE("/:id", jwt.JWTAuthMiddleware(), server.DeleteUserHandler)
	}
}
