package controllers

import (
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
	// headerのtokenを取得
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Error invalid authorization token",
		})
		return
	}

	issuesReq := types.IssuesReq{}

	if err := c.ShouldBindQuery(&issuesReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	i, err := ic.m.GetIssues(token, issuesReq)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, i)
}

func (ic *issueController) Get(c *gin.Context) {
	// headerのtokenを取得
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Error invalid authorization token",
		})
		return
	}

	issueReq := types.IssueReq{}

	if err := c.ShouldBindQuery(&issueReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	i, err := ic.m.GetIssue(token, issueReq)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, i)
}

func (ic *issueController) Create(c *gin.Context) {
	// headerのtokenを取得
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Error invalid authorization token",
		})
		return
	}

	issueCreateReq := types.IssueCreateReq{}

	if err := c.ShouldBindJSON(&issueCreateReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	i, err := ic.m.CreateIssue(token, issueCreateReq)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, i)
}
