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
			assert.NotEqual(t, v.ID, "")
			assert.NotEqual(t, v.CreatedAt, "")
			assert.NotEqual(t, v.UpdatedAt, "")
			assert.NotEqual(t, v.URL, "")
			assert.NotEqual(t, v.State, "")
			assert.NotEqual(t, v.Title, "")
			assert.NotEqual(t, v.Number, "")
			assert.NotEqual(t, v.Body, "")
			assert.NotEqual(t, v.BodyHTML, "")
		}
	})
}
