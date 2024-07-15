package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/lahnasti/GO_praktikum/internal/server"
)

func UserRoutes(r *gin.Engine, server *server.Server) {
	userGroup := r.Group("/users")
	{
		userGroup.GET("/", server.GetUsersHandler)
		userGroup.POST("/add", server.RegisterUserHandler)
		userGroup.POST("/adds", server.RegisterMultipleUsersHadler)
		userGroup.POST("/login", server.LoginHandler)

		userGroup.GET("/:id", server.JWTAuthMiddleware(), server.GetUserByIDHandler)
		userGroup.PUT("/:id", server.JWTAuthMiddleware(), server.UpdateUserHandler)
		userGroup.DELETE("/:id", server.JWTAuthMiddleware(), server.DeleteUserHandler)
	}
}

func BookRoutes(r *gin.Engine, server *server.Server) {
	bookGroup := r.Group("/books")
	{
		bookGroup.GET("/", server.GetBooksHandler)
		bookGroup.POST("/add", server.JWTAuthMiddleware(), server.CreateBookHandler)
		bookGroup.POST("/adds", server.JWTAuthMiddleware(), server.CreateMultipleBooksHandler)

		bookGroup.GET("/:user_id", server.JWTAuthMiddleware(), server.GetBooksByUserHandler)
	}
}
