package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/tangoctrung/golang_api_v2/dto"
	"github.com/tangoctrung/golang_api_v2/entity"
	"github.com/tangoctrung/golang_api_v2/helper"
	"github.com/tangoctrung/golang_api_v2/service"
)

type BookController interface {
	InsertBook(ctx *gin.Context)
	UpdateBook(ctx *gin.Context)
	DeleteBook(ctx *gin.Context)
	FindBookByID(ctx *gin.Context)
	GetAllBooks(ctx *gin.Context)
}

type bookController struct {
	bookService service.BookService
	jwtService  service.JWTService
}

func NewBookController(bookService service.BookService, jwtService service.JWTService) BookController {
	return &bookController{
		bookService: bookService,
		jwtService:  jwtService,
	}
}

func (c *bookController) InsertBook(ctx *gin.Context) {
	var bookDTO dto.BookCreateDTO
	errDTO := ctx.ShouldBind(&bookDTO)
	if errDTO != nil {
		res := helper.BuildErrorsResponse(false, "Failed to create book", errDTO.Error())
		ctx.JSON(http.StatusBadRequest, res)
	}

	authHeader := ctx.GetHeader("Authorization")
	userID := c.getUserIDByToken(authHeader)
	convertedUserID, err := strconv.ParseUint(userID, 10, 64)
	if err == nil {
		bookDTO.UserID = convertedUserID
	}
	result := c.bookService.InsertBook(bookDTO)
	response := helper.BuildResponse(true, "OK", result)
	ctx.JSON(http.StatusCreated, response)
}

func (c *bookController) UpdateBook(ctx *gin.Context) {
	var bookUpdateDTO dto.BookUpdateDTO
	errDTO := ctx.ShouldBind(&bookUpdateDTO)
	if errDTO != nil {
		res := helper.BuildErrorsResponse(false, "Failed to process request", errDTO.Error())
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	authHeader := ctx.GetHeader("Authorization")
	token, errToken := c.jwtService.ValidateToken(authHeader)
	if errToken != nil {
		panic(errToken.Error())
	}
	claims := token.Claims.(jwt.MapClaims)
	userID := fmt.Sprintf("%v", claims["user_id"])
	if c.bookService.IsAllowedToEdit(userID, bookUpdateDTO.ID) {
		id, errID := strconv.ParseUint(userID, 10, 64)
		if errID == nil {
			bookUpdateDTO.UserID = id
		}
		result := c.bookService.UpdateBook(bookUpdateDTO)
		response := helper.BuildResponse(true, "OK", result)
		ctx.JSON(http.StatusOK, response)
	} else {
		response := helper.BuildErrorsResponse(false, "You dont have permission", "You are not the owner")
		ctx.JSON(http.StatusForbidden, response)
	}

}

func (c *bookController) DeleteBook(ctx *gin.Context) {
	var book entity.Book
	id, err := strconv.ParseUint(ctx.Param("id"), 0, 0)
	if err != nil {
		response := helper.BuildErrorsResponse(false, "Failed tou get id", "No param id were found")
		ctx.JSON(http.StatusBadRequest, response)
	}
	book.ID = id
	authHeader := ctx.GetHeader("Authorization")
	token, errToken := c.jwtService.ValidateToken(authHeader)
	if errToken != nil {
		panic(errToken.Error())
	}
	claims := token.Claims.(jwt.MapClaims)
	userID := fmt.Sprintf("%v", claims["user_id"])
	if c.bookService.IsAllowedToEdit(userID, book.ID) {
		c.bookService.DeleteBook(book)
		res := helper.BuildResponse(true, "Deleted", nil)
		ctx.JSON(http.StatusOK, res)
	} else {
		response := helper.BuildErrorsResponse(false, "You dont have permission", "You are not the owner")
		ctx.JSON(http.StatusForbidden, response)
	}

}

func (c *bookController) FindBookByID(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 0, 0)
	if err != nil {
		res := helper.BuildErrorsResponse(false, "No param id was found", err.Error())
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	var book entity.Book = c.bookService.FindBookByID(id)
	if (book == entity.Book{}) {
		res := helper.BuildErrorsResponse(false, "Data not found", "No data with given id")
		ctx.JSON(http.StatusNotFound, res)
	} else {
		res := helper.BuildResponse(true, "OK", book)
		ctx.JSON(http.StatusOK, res)
	}

}

func (c *bookController) GetAllBooks(ctx *gin.Context) {
	var books []entity.Book = c.bookService.GetAllBooks()
	res := helper.BuildResponse(true, "OK", books)
	ctx.JSON(http.StatusOK, res)

}

func (c *bookController) getUserIDByToken(token string) string {
	aToken, err := c.jwtService.ValidateToken(token)
	if err != nil {
		panic(err.Error())
	}
	claims := aToken.Claims.(jwt.MapClaims)
	id := fmt.Sprintf("%v", claims["user_id"])
	return id
}
