package middleware

import (
	"main/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ParseAuthHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Error invalid authorization token",
			})
			c.Abort()
			return
		}

		token, err := utils.RefreshJWT("token", token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			c.Abort()
			return
		}

		c.Set("token", token)
		// response headerには、jwtをセットする
		jwtStr, err := utils.GenerateJWT("token", token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			c.Abort()
			return
		}
		c.Header("Authorization", jwtStr)
	}
}
