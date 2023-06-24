package models

import (
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
		owner := os.Getenv("OWNER_NAME")   // オーナー名
		repo := os.Getenv("REPO_NAME")     // リポジトリ名
		first := 5                         // 取得数
		order := "DESC"                    // 取得順
		states := "OPEN"                   // OPENかCLOSE

		model := NewIssueModel()

		i, err := model.GetIssues(token, owner, repo, first, order, states)
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
