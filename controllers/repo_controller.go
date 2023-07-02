package controllers

import (
	"main/models"
	"main/types"
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

func (rc *repoController) IndexByToken(c *gin.Context) {
	// headerのtokenを取得
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Error invalid authorization token",
		})
		return
	}

	reposReq := types.ReposReq{}

	if err := c.ShouldBindQuery(&reposReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	r, err := rc.m.GetReposByToken(token, reposReq)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, r)
}
