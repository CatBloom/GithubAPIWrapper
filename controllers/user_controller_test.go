package controllers

import (
	"errors"
	"main/types"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type mockUserModel struct{}

func (m *mockUserModel) GetUserByToken(_ string) (types.UserRes, error) {
	user := types.User{
		Login:     "john_doe",
		Name:      "John Doe",
		URL:       "https://example.com/john_doe",
		AvatarUrl: "https://example.com/avatar/john_doe.jpg",
	}

	res := types.UserRes{
		Data: types.UserViewer{
			Viewer: user,
		},
	}

	return res, nil
}

type mockErrorUserModel struct{}

func (m *mockErrorUserModel) GetUserByToken(_ string) (types.UserRes, error) {
	return types.UserRes{}, errors.New("Failed to user")
}

func TestGetByTokenUserController(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Run("Test with success response", func(t *testing.T) {
		mockModel := &mockUserModel{}
		controller := NewUserController(mockModel)

		res := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(res)
		c.Request, _ = http.NewRequest(http.MethodGet, "/user", nil)
		c.Set("token", "token")
		c.Request.Header.Set("Authorization", "Bearer")

		controller.GetByToken(c)

		assert.Equal(t, http.StatusOK, res.Code)
		expectedResponse := `{"data":{"viewer":{"login":"john_doe","name":"John Doe","url":"https://example.com/john_doe","avatarUrl":"https://example.com/avatar/john_doe.jpg"}}}`
		assert.Equal(t, expectedResponse, res.Body.String())
	})

	t.Run("Test with error response", func(t *testing.T) {
		mockErrorModel := &mockErrorUserModel{}
		controller := NewUserController(mockErrorModel)

		res := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(res)
		c.Request, _ = http.NewRequest(http.MethodGet, "/user", nil)
		c.Set("token", "token")
		c.Request.Header.Set("Authorization", "Bearer")

		controller.GetByToken(c)

		assert.Equal(t, http.StatusNotFound, res.Code)
		expectedResponse := `{"error":"Failed to user"}`
		assert.Equal(t, expectedResponse, res.Body.String())
	})
}
