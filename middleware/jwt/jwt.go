package jwt

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"

	"github.com/EDDYCJY/go-gin-example/pkg/util"
)

// JWT is jwt middleware
func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		// var code int
		// var data interface{}

		// code = e.SUCCESS
		token := c.Query("token")
		if token == "" {
			// code = e.INVALID_PARAMS
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "token is empty",
			})
			// Abort 确保中间件链中的其他处理程序不会被调用
			c.Abort()
			return

		} else {
			_, err := util.ParseToken(token)
			if err != nil {
				switch err.(*jwt.ValidationError).Errors {
				case jwt.ValidationErrorExpired:
					// code = e.ERROR_AUTH_CHECK_TOKEN_TIMEOUT

					c.JSON(http.StatusUnauthorized, gin.H{
						"error": "token is expired",
					})
					c.Abort()
					return

				default:
					// code = e.ERROR_AUTH_CHECK_TOKEN_FAIL
					c.JSON(http.StatusUnauthorized, gin.H{
						"error": "token is invalid",
					})
					c.Abort()
					return
				}
			}
		}

		c.Next()
	}
}
