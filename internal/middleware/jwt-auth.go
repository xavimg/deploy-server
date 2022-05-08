package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xavimg/Turing/apituringserver/internal/helper"
	"github.com/xavimg/Turing/apituringserver/internal/service"
)

// AuthorizeJWT validates the token user given, return 401 if not valid
func AuthorizeJWT(jwtService service.JWTService) gin.HandlerFunc {
	return func(context *gin.Context) {
		authHeader := context.GetHeader("Authorization")
		googleLogin := context.Param("state")

		if googleLogin == "randomstate" {
			context.Next()
		}

		if authHeader == "" {
			response := helper.BuildErrorResponse("Failed to process request", "No token found", nil)
			context.AbortWithStatusJSON(http.StatusBadRequest, response)
			return
		}

		token, err := jwtService.ValidateToken(authHeader)

		if !token.Valid {
			response := helper.BuildErrorResponse("Token is not valid", err.Error(), nil)
			context.AbortWithStatusJSON(http.StatusUnauthorized, response)
		}

		context.Next()
	}
}
