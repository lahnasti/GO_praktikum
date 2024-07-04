package routes

import 	(
	"github.com/gin-gonic/gin"
	"github.com/lahnasti/GO_praktikum/internal/server/books"
)

func BookRoutes(r *gin.Engine, server *books.Server) {
	bookGroup := r.Group("/books")
	{
		bookGroup.GET("/", server.GetBooksHandler)
		bookGroup.POST("/", server.CreateBookHandler)

	}
}