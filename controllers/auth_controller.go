package controllers

import (
	"main/models"
	"main/types"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type AuthController interface {
	Get(c *gin.Context)
}

type authController struct {
	m models.AuthModel
}

func NewAuthController(m models.AuthModel) AuthController {
	return &authController{m}
}

func (ac *authController) Get(c *gin.Context) {
	authReq := types.AuthReq{}

	if err := c.ShouldBindQuery(&authReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	t, err := ac.m.GetToken(authReq)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	cookie := &http.Cookie{
		Name:     "Token",
		Value:    cases.Title(language.Und).String(t.TokenType) + " " + t.AccessToken,
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
	}

	http.SetCookie(c.Writer, cookie)

	c.JSON(http.StatusOK, "success get token")
}
