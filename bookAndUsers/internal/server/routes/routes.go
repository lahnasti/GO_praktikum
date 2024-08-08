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

		userGroup.GET("/:id", server.GetUserHandler)
		userGroup.PUT("/:id", server.UpdateUserHandler)
		userGroup.DELETE("/:id", server.DeleteUserHandler)
	}
}

func BookRoutes(r *gin.Engine, server *server.Server) {
	bookGroup := r.Group("/books")
	{
		bookGroup.GET("/", server.GetAllBooksHandler)
		bookGroup.POST("/add", server.SaveBookHandler)
		bookGroup.POST("/adds", server.SaveBooksHandler)

		bookGroup.GET("/:user_id", server.GetBooksByUserHandler)
		bookGroup.DELETE("/delete/:id", server.DeleteBookHandler)
	}
}
