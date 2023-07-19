package controllers

import (
	"fmt"
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
	token, exists := c.Get("token")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Error invalid token",
		})
		return
	}
	sToken := fmt.Sprint(token)

	reposReq := types.ReposReq{}

	if err := c.ShouldBindQuery(&reposReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	r, err := rc.m.GetReposByToken(sToken, reposReq)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, r)
}
