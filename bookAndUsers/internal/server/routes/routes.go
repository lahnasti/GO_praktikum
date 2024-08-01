package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/lahnasti/GO_praktikum/internal/server"
)

func UserRoutes(r *gin.Engine, server *server.Server) {
	userGroup := r.Group("/users")
	{
		userGroup.GET("/", server.GetAllUsersHandler)
		userGroup.POST("/register", server.RegisterUserHandler)
		userGroup.POST("/login", server.LoginHandler)

		userGroup.GET("/:id", server.JWTAuthMiddleware(), server.GetUserByIDHandler)
		userGroup.PUT("/:id", server.JWTAuthMiddleware(), server.UpdateUserHandler)
		userGroup.DELETE("/:id", server.JWTAuthMiddleware(), server.DeleteUserHandler)
	}
}

func BookRoutes(r *gin.Engine, server *server.Server) {
	bookGroup := r.Group("/books")
	{
		bookGroup.GET("/", server.GetAllBooksHandler)
		bookGroup.POST("/add", server.JWTAuthMiddleware(), server.SaveBookHandler)
		bookGroup.POST("/adds", server.JWTAuthMiddleware(), server.SaveBooksHandler)

		bookGroup.GET("/:user_id", server.JWTAuthMiddleware(), server.GetBooksByUserHandler)
		bookGroup.DELETE("/delete/:id", server.JWTAuthMiddleware(), server.DeleteBookHandler)	}
}
