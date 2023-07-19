package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ParseAuthHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Cookieを取得
		token, err := c.Cookie("Token")

		// ParseAuthHandlerでは、cookieからtokenを取得できない場合
		// http.StatusAccepted(202)でmessageを返す
		if err != nil {
			c.JSON(http.StatusAccepted, gin.H{
				"message": "token is missing in cookie",
			})
			c.Abort()
			return
		}
		if token == "" {
			c.JSON(http.StatusAccepted, gin.H{
				"message": "token is empty in cookie",
			})
			c.Abort()
			return
		}

		c.Set("token", token)
	}
}
