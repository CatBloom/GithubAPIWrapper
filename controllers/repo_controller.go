package controllers

import (
	"main/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RepoController interface {
	IndexByToken(c *gin.Context)
}

type repoController struct {
	m models.RepoModel
}

func NewRepoController(m models.RepoModel) RepoController {
	return &repoController{m}
}

type RepoReq struct {
	First int    `form:"first" binding:"required,max=100,min=1"`
	Order string `form:"order" binding:"required,oneof=ASC DESC"`
}

func (rc repoController) IndexByToken(c *gin.Context) {
	// headerのtokenを取得
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Error invalid authorization token",
		})
		return
	}

	repoReq := RepoReq{}

	if err := c.ShouldBindQuery(&repoReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	r, err := rc.m.GetReposByToken(token, repoReq.First, repoReq.Order)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, r)
}
