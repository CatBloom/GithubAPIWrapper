package models

import (
	"main/utils"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserModels(t *testing.T) {
	t.Logf("models/user test")

	utils.InitEnv()
	t.Run("Test with success response", func(t *testing.T) {
		// テスト用に環境変数からアクセストークンを取得
		token := os.Getenv("ACCESS_TOKEN")
		models := NewUserModels()

		u, err := models.GetUser(token)
		if err != nil {
			t.Errorf("GetUser returned an error: %v", err)
		}

		assert.NotEqual(t, u.Data.Login, "")
		assert.NotEqual(t, u.Data.Name, "")
		assert.NotEqual(t, u.Data.Url, "")
		assert.NotEqual(t, u.Data.AvatarUrl, "")
	})
}
