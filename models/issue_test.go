package models

import (
	"main/types"
	"main/utils"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetIssuesIssueModel(t *testing.T) {
	gin.SetMode(gin.TestMode)
	utils.InitEnv()
	t.Run("Test with success response", func(t *testing.T) {
		// テスト用に環境変数から各項目を取得
		token := os.Getenv("ACCESS_TOKEN") // アクセストークン
		issuesReq := types.IssuesReq{
			Owner:  os.Getenv("OWNER_NAME"), // オーナー名
			Repo:   os.Getenv("REPO_NAME"),  // リポジトリ名
			First:  5,                       // 取得数
			Order:  "DESC",                  // 取得順
			States: "OPEN",                  // OPENかCLOSE
			After:  "",                      // 次ページ取得
		}

		model := NewIssueModel()

		i, err := model.GetIssues(token, issuesReq)
		if err != nil {
			t.Errorf("GetIssues returned an error: %v", err)
		}

		for _, v := range i.Data.Repository.Issues.Nodes {
			assert.NotEmpty(t, v.ID)
			assert.NotEmpty(t, v.CreatedAt)
			assert.NotEmpty(t, v.UpdatedAt)
			assert.NotEmpty(t, v.URL)
			assert.NotEmpty(t, v.State)
			assert.NotEmpty(t, v.Title)
			assert.NotEmpty(t, v.Number)
			assert.NotEmpty(t, v.Body)
			assert.NotEmpty(t, v.BodyHTML)
		}
	})
}

func TestGetIssueIssueModel(t *testing.T) {
	gin.SetMode(gin.TestMode)
	utils.InitEnv()
	t.Run("Test with success response", func(t *testing.T) {
		// テスト用に環境変数から各項目を取得
		token := os.Getenv("ACCESS_TOKEN") // アクセストークン
		issueReq := types.IssueReq{
			Owner:  os.Getenv("OWNER_NAME"), // オーナー名
			Repo:   os.Getenv("REPO_NAME"),  // リポジトリ名
			Number: 1,
		}

		model := NewIssueModel()

		i, err := model.GetIssue(token, issueReq)
		if err != nil {
			t.Errorf("GetIssues returned an error: %v", err)
		}

		v := i.Data.Repository.Issue
		assert.NotEmpty(t, v.ID)
		assert.NotEmpty(t, v.CreatedAt)
		assert.NotEmpty(t, v.UpdatedAt)
		assert.NotEmpty(t, v.URL)
		assert.NotEmpty(t, v.State)
		assert.NotEmpty(t, v.Title)
		assert.NotEmpty(t, v.Number)
		assert.NotEmpty(t, v.Body)
		assert.NotEmpty(t, v.BodyHTML)
	})
}
