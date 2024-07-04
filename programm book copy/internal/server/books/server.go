package books

import (
	"net/http"
	//"log"
	//"strconv"
	//"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/lahnasti/GO_praktikum/internal/models"
	//"github.com/rs/zerolog"
)

type Repository interface {
	CreateBook(models.Book) (string, error)
	GetBooks() ([]models.Book, error)
	//GetBookByID(string)(models.Book, error)
}

type Server struct {
	Db Repository
	Valid *validator.Validate
	//log *zerolog.Logger
}

func (s *Server) GetBooksHandler(ctx *gin.Context) {
	books, err := s.Db.GetBooks()
	if err != nil {
		//s.log.Error().Err(err).Msg("Failed inquiry")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{"message": "Storage books", "books": books})
}

func (s *Server) CreateBookHandler(ctx *gin.Context) {
	var book models.Book
	err := ctx.ShouldBindBodyWithJSON(&book)
	if err != nil {
		//s.log.Error().Err(err).Msg("Failed unmarshal body")
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid params", "error": err.Error()})
		return
	}
	err = s.Valid.Struct(book)
	if err != nil {
		//s.log.Error().Err(err).Msg("Failed validation")
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Data has not been validated", "error": err.Error()})
		return
	}
	bookID, err := s.Db.CreateBook(book)
	if err != nil {
		// if strings.Contains(err.Error(), "permission denied for table users") {
          // ctx.JSON(http.StatusForbidden, gin.H{"message": "Permission denied for table users"})
		//} else {
		//s.log.Error().Err(err).Msg("Failed to save book")
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to save book"})
		return
	}
	ctx.JSON(200, gin.H{"message": "Book successfully created", "book_id": bookID})

	}

	/*func (s *Server) GetBookByIDHandler(ctx *gin.Context) {
		param := ctx.Query("id")
		log.Println("Param from url - " + param)
		id, err := strconv.Atoi(param)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid argument", "error": err.Error()})
			return
		}
		log.Printf("ID - %v", id)
		book, err := s.Db.GetBookByID(id)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(200, gin.H{"message": "Book retrieved", "book": book})
	}*/
