package models

import (
	"main/utils"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestGetTokenAuthModel(t *testing.T) {
	gin.SetMode(gin.TestMode)
	utils.InitEnv()

	// GetTokenは一時コードが必要なため、テストは行わない
	// 雛形だけ作成
	t.Run("Test with success response", func(t *testing.T) {
		// authReq := types.AuthReq{
		// 	Code: "abcd12345", //一時コード
		// }

		// model := NewAuthModel()

		// _, err := model.GetToken(authReq)
		// if err != nil {
		// 	t.Errorf("GetIssues returned an error: %v", err)
		// }
	})
}
