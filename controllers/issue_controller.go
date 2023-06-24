package controllers

import (
	"main/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type IssueController interface {
	Index(c *gin.Context)
}

type issueController struct {
	m models.IssueModel
}

func NewIssueController(m models.IssueModel) IssueController {
	return &issueController{m}
}

type IssueReq struct {
	Owner  string `form:"owner" binding:"required"`
	Repo   string `form:"repo" binding:"required"`
	First  int    `form:"first" binding:"required,max=100,min=1"`
	Order  string `form:"order" binding:"required,oneof=ASC DESC"`
	States string `form:"states" binding:"required,oneof=OPEN CLOSE"`
}

func (ic issueController) Index(c *gin.Context) {
	// headerのtokenを取得
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Error invalid authorization token",
		})
		return
	}

	issueReq := IssueReq{}

	if err := c.ShouldBindQuery(&issueReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	i, err := ic.m.GetIssues(token, issueReq.Owner, issueReq.Repo, issueReq.First, issueReq.Order, issueReq.States)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, i)
}
