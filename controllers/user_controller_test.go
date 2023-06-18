package controllers

import (
	"errors"
	"main/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type mockUserModels struct{}

func (m *mockUserModels) GetUser(_ string) (models.UserResp, error) {
	user := models.User{
		Login:     "john_doe",
		Name:      "John Doe",
		URL:       "https://example.com/john_doe",
		AvatarUrl: "https://example.com/avatar/john_doe.jpg",
	}

	resp := models.UserResp{
		Data: models.UserViewer{
			Viewer: user,
		},
	}

	userResp := resp
	return userResp, nil
}

type mockErrorUserModels struct{}

func (m *mockErrorUserModels) GetUser(_ string) (models.UserResp, error) {
	errorMessage := "Failed to user"
	return models.UserResp{}, errors.New(errorMessage)
}

func TestGetUserControllers(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Run("Test with success response", func(t *testing.T) {
		mockModels := &mockUserModels{}
		controller := NewUserController(mockModels)

		res := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(res)
		c.Request, _ = http.NewRequest(http.MethodGet, "/user", nil)
		// c.Set("token", "token")
		c.Request.Header.Set("Authorization", "Bearer")

		controller.Get(c)

		assert.Equal(t, http.StatusOK, res.Code)
		expectedResponse := `{"data":{"viewer":{"login":"john_doe","name":"John Doe","url":"https://example.com/john_doe","avatarUrl":"https://example.com/avatar/john_doe.jpg"}}}`
		assert.Equal(t, expectedResponse, res.Body.String())
	})

	t.Run("Test with error response", func(t *testing.T) {
		mockErrorModels := &mockErrorUserModels{}
		controller := NewUserController(mockErrorModels)

		res := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(res)
		c.Request, _ = http.NewRequest(http.MethodGet, "/user", nil)
		// c.Set("token", "token")
		c.Request.Header.Set("Authorization", "Bearer")

		controller.Get(c)

		assert.Equal(t, http.StatusNotFound, res.Code)
		expectedResponse := `{"error":"Failed to user"}`
		assert.Equal(t, expectedResponse, res.Body.String())
	})
}
