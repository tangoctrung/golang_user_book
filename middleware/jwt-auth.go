package middleware

import (
	"log"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/tangoctrung/golang_api_v2/helper"
	"github.com/tangoctrung/golang_api_v2/service"
)

func AuthorizeJWT(jwtService service.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response := helper.BuildErrorsResponse(false, "Failed to process request", "No token found")
			c.AbortWithStatusJSON(http.StatusBadRequest, response)
			return
		}
		token, err := jwtService.ValidateToken(authHeader)
		if token.Valid {
			claims := token.Claims.(jwt.MapClaims)
			log.Println("Claims[user_id]: ", claims["user_id"])
			log.Println("Claims[issuer]: ", claims["issuer"])
		} else {
			log.Println(err)
			response := helper.BuildErrorsResponse(false, "Token is valid", err.Error())
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
		}
	}
}
