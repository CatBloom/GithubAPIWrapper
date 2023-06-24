package models

import (
	"main/utils"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetReposByTokenRepoModel(t *testing.T) {
	gin.SetMode(gin.TestMode)
	utils.InitEnv()
	t.Run("Test with success response", func(t *testing.T) {
		// テスト用に環境変数からアクセストークンを取得
		token := os.Getenv("ACCESS_TOKEN")
		first := 5      // 取得数
		order := "DESC" // 取得順

		model := NewRepoModel()

		r, err := model.GetReposByToken(token, first, order)
		if err != nil {
			t.Errorf("GetReposByToken returned an error: %v", err)
		}

		for _, v := range r.Data.Viewer.Repositories.Nodes {
			assert.NotEmpty(t, v.Name)
			assert.NotEmpty(t, v.URL)
			assert.NotEmpty(t, v.CreatedAt)
			assert.NotEmpty(t, v.UpdatedAt)
		}
	})
}
