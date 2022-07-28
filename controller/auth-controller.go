package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tangoctrung/golang_api_v2/dto"
	"github.com/tangoctrung/golang_api_v2/entity"
	"github.com/tangoctrung/golang_api_v2/helper"
	"github.com/tangoctrung/golang_api_v2/service"
)

type AuthController interface {
	Login(ctx *gin.Context)
	Register(ctx *gin.Context)
}

type authController struct {
	authService service.AuthService
	jwtService  service.JWTService
}

func NewAuthController(authService service.AuthService, jwtService service.JWTService) AuthController {
	return &authController{
		authService: authService,
		jwtService:  jwtService,
	}
}

func (c *authController) Login(ctx *gin.Context) {
	var loginDTO dto.LoginDTO
	errDTO := ctx.ShouldBind(&loginDTO)
	if errDTO != nil {
		reponse := helper.BuildErrorsResponse(false, "Failed to process login request", errDTO.Error())
		ctx.AbortWithStatusJSON(http.StatusBadRequest, reponse)
	}
	authResult := c.authService.VerifyCredential(loginDTO.Email, loginDTO.Password)
	if v, ok := authResult.(entity.User); ok {
		generatedToken := c.jwtService.GenerateToken(strconv.FormatUint(v.ID, 10))
		v.Token = generatedToken
		response := helper.BuildResponse(true, "ok", v)
		ctx.JSON(http.StatusOK, response)
		return
	}
	response := helper.BuildErrorsResponse(false, "", "Invalid Credential")
	ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
}

func (c *authController) Register(ctx *gin.Context) {
	var registerDTO dto.RegisterDTO
	errDTO := ctx.ShouldBind(&registerDTO)
	if errDTO != nil {
		response := helper.BuildErrorsResponse(false, "Failed to process register request", errDTO.Error())
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}
	if !c.authService.IsDuplicateEmail(registerDTO.Email) {
		response := helper.BuildErrorsResponse(false, "Failed to process register request", "Email already registered")
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
	} else {
		createdUser := c.authService.CreateUser(registerDTO)
		token := c.jwtService.GenerateToken(strconv.FormatUint(createdUser.ID, 10))
		createdUser.Token = token
		response := helper.BuildResponse(true, "OK!", createdUser)
		ctx.AbortWithStatusJSON(http.StatusCreated, response)
	}
}
