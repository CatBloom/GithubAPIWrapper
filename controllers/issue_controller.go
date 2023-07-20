package controllers

import (
	"fmt"
	"main/models"
	"main/types"
	"net/http"

	"github.com/gin-gonic/gin"
)

type IssueController interface {
	Index(c *gin.Context)
	Get(c *gin.Context)
	Create(c *gin.Context)
}

type issueController struct {
	m models.IssueModel
}

func NewIssueController(m models.IssueModel) IssueController {
	return &issueController{m}
}

func (ic *issueController) Index(c *gin.Context) {
	token, exists := c.Get("token")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Error invalid token",
		})
		return
	}
	sToken := fmt.Sprint(token)

	issuesReq := types.IssuesReq{}

	if err := c.ShouldBindQuery(&issuesReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	i, err := ic.m.GetIssues(sToken, issuesReq)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, i)
}

func (ic *issueController) Get(c *gin.Context) {
	token, exists := c.Get("token")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Error invalid token",
		})
		return
	}
	sToken := fmt.Sprint(token)

	issueReq := types.IssueReq{}

	if err := c.ShouldBindQuery(&issueReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	i, err := ic.m.GetIssue(sToken, issueReq)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, i)
}

func (ic *issueController) Create(c *gin.Context) {
	token, exists := c.Get("token")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Error invalid token",
		})
		return
	}
	sToken := fmt.Sprint(token)

	issueCreateReq := types.IssueCreateReq{}

	if err := c.ShouldBindJSON(&issueCreateReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	i, err := ic.m.CreateIssue(sToken, issueCreateReq)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, i)
}
