package models

import (
	"main/utils"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetUserByTokenUserModel(t *testing.T) {
	gin.SetMode(gin.TestMode)
	utils.InitEnv()
	t.Run("Test with success response", func(t *testing.T) {
		// テスト用に環境変数からアクセストークンを取得
		token := os.Getenv("ACCESS_TOKEN")
		model := NewUserModel()

		u, err := model.GetUserByToken(token)
		if err != nil {
			t.Errorf("GetUser returned an error: %v", err)
		}

		assert.NotEmpty(t, u.Data.Viewer.Login)
		assert.NotEmpty(t, u.Data.Viewer.Name)
		assert.NotEmpty(t, u.Data.Viewer.URL)
		assert.NotEmpty(t, u.Data.Viewer.AvatarUrl)
	})
}
