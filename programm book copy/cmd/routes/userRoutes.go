package routes

import 	(
	"github.com/gin-gonic/gin"
	"github.com/lahnasti/GO_praktikum/internal/server/users"
	"github.com/lahnasti/GO_praktikum/internal/server/users/jwt"
)

func UserRoutes(r *gin.Engine, server *users.Server) {
	userGroup := r.Group("/users")
	{
		userGroup.GET("/", server.GetUsersHandler)
		userGroup.POST("/add", jwt.JWTAuthMiddleware(), server.RegisterUserHandler)
		userGroup.POST("/adds", jwt.JWTAuthMiddleware(), server.RegisterMultipleUsersHadler)

		userGroup.GET("/:id", server.GetUserByIDHandler)
		userGroup.PUT("/:id", server.UpdateUserHandler)
		userGroup.DELETE("/:id", server.DeleteUserHandler)
	}
}

