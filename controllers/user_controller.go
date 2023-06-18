package controllers

import (
	"main/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserController interface {
	GetByToken(c *gin.Context)
}

type userController struct {
	m models.UserModels
}

func NewUserController(m models.UserModels) UserController {
	return &userController{m}
}

func (uc userController) GetByToken(c *gin.Context) {
	// auth_middlewareを使用する際の処理
	// token, exists := c.Get("token")
	// if !exists {
	// 	c.JSON(http.StatusUnauthorized, gin.H{
	// 		"error": "Error invalid token",
	// 	})
	// 	return
	// }
	// sToken := fmt.Sprint(token)

	// headerのtokenを取得
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Error invalid authorization token",
		})
		return
	}

	u, err := uc.m.GetUserByToken(token)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, u)
}
