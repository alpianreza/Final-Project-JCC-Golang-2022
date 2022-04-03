package middlewares

import (
	"finalproject/utils/token"
	"net/http"

	"github.com/gin-gonic/gin"
)

func JwtAuthMiddleware(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := token.TokenValid(c)
		if err != nil {
			c.String(http.StatusUnauthorized, err.Error())
			c.Abort()
			return
		}
		roleToken, err := token.ExtractTokenRole(c)
		if err != nil {
			c.String(http.StatusUnauthorized, err.Error())
			c.Abort()
			return
		}
		for _, role := range roles {
			if role == roleToken {
				c.Next()
				return
			}
		}
		c.String(http.StatusUnauthorized, "Role not match")
		c.Abort()
	}
}
