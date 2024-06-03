package middleware

import (
	"net/http"
	"strings"

	"github.com/AhmadSafrizal/myGram/helper"
	"github.com/gin-gonic/gin"
)

func Authotization() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		headerAuth := ctx.GetHeader("Authorization")

		splitToken := strings.Split(headerAuth, " ")

		if len(splitToken) != 2 {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, map[string]any{
				"message": "invalid authorization header",
			})
			return
		}

		if splitToken[0] != "Bearer" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, map[string]any{
				"message": "invalid authorization method",
			})
			return
		}

		tokenValid := helper.ValidateUserJWT(splitToken[1])

		if !tokenValid {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, map[string]any{
				"message": "malformed token",
			})
			return
		}

		ctx.Next()
	}
}
