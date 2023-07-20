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

type mockAuthModel struct{}

func (m *mockAuthModel) GetToken(types.AuthReq) (types.AuthRes, error) {
	res := types.AuthRes{
		AccessToken: "abcde12345",
		TokenType:   "Bearer",
		Scope:       "repo",
	}

	return res, nil
}

type mockErrorAuthModel struct{}

func (m *mockErrorAuthModel) GetToken(types.AuthReq) (types.AuthRes, error) {
	return types.AuthRes{}, errors.New("Failed to token")
}

func TestGetAuthController(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Run("Test with success response", func(t *testing.T) {
		mockModel := &mockAuthModel{}
		controller := NewAuthController(mockModel)

		res := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(res)
		c.Request, _ = http.NewRequest(http.MethodGet, "/auth?code=abcdefg12345", nil)

		controller.Get(c)

		assert.Equal(t, http.StatusOK, res.Code)
	})

	t.Run("Test with error response", func(t *testing.T) {
		mockErrorModel := &mockErrorAuthModel{}
		controller := NewAuthController(mockErrorModel)

		res := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(res)
		c.Request, _ = http.NewRequest(http.MethodGet, "/auth?code=abcdefg12345", nil)

		controller.Get(c)

		assert.Equal(t, http.StatusNotFound, res.Code)
		expectedResponse := `{"error":"Failed to token"}`
		assert.Equal(t, expectedResponse, res.Body.String())
	})

	// QueryParam
	paramCases := []struct {
		name   string
		params string
	}{
		{
			name:   "validate error handling for query code",
			params: "code=",
		},
	}

	for _, tc := range paramCases {
		t.Run("Test "+tc.name, func(t *testing.T) {
			mockModel := &mockAuthModel{}
			controller := NewAuthController(mockModel)

			res := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(res)
			url := "/auth?" + tc.params
			c.Request, _ = http.NewRequest(http.MethodGet, url, nil)

			controller.Get(c)

			assert.Equal(t, http.StatusBadRequest, res.Code)
		})
	}
}
