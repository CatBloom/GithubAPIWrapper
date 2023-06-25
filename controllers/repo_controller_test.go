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

type mockRepoModel struct{}

func (m *mockRepoModel) GetReposByToken(_ string, first int, order string, after string) (models.RepoResp, error) {
	repo := []models.Repo{
		{
			Name:        "First Repository",
			URL:         "https://example.com/repository",
			Description: "description",
			CreatedAt:   "2000-01-01T00:00:00Z",
			UpdatedAt:   "2000-01-01T00:00:00Z",
		},
		{
			Name:        "Second Repository",
			URL:         "https://example.com/repository",
			Description: "",
			CreatedAt:   "2000-01-01T00:00:00Z",
			UpdatedAt:   "2000-01-01T00:00:00Z",
		},
	}

	resp := models.RepoResp{
		Data: models.RepoViewer{
			Viewer: models.RepoRepositories{
				Repositories: models.RepoNodes{
					Nodes: repo,
					PageInfo: models.RepoPageInfo{
						EndCutsor:   "",
						HasNextPage: false,
					},
				},
			},
		},
	}

	repoResp := resp
	return repoResp, nil
}

type mockErrorRepoModel struct{}

func (m *mockErrorRepoModel) GetReposByToken(_ string, first int, order string, after string) (models.RepoResp, error) {
	errorMessage := "Failed to repositories"
	return models.RepoResp{}, errors.New(errorMessage)
}

func TestIndexByTokenRepoController(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Run("Test with success response", func(t *testing.T) {
		mockModel := &mockRepoModel{}
		controller := NewRepoController(mockModel)

		res := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(res)
		c.Request, _ = http.NewRequest(http.MethodGet, "/repo?first=1&order=DESC", nil)
		c.Set("token", "token")
		c.Request.Header.Set("Authorization", "Bearer")

		controller.IndexByToken(c)

		assert.Equal(t, http.StatusOK, res.Code)
	})

	t.Run("Test with error response", func(t *testing.T) {
		mockErrorModel := &mockErrorRepoModel{}
		controller := NewRepoController(mockErrorModel)

		res := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(res)
		c.Request, _ = http.NewRequest(http.MethodGet, "/repo?first=1&order=DESC", nil)
		c.Request.Header.Set("Authorization", "Bearer")

		controller.IndexByToken(c)

		assert.Equal(t, http.StatusNotFound, res.Code)
		expectedResponse := `{"error":"Failed to repositories"}`
		assert.Equal(t, expectedResponse, res.Body.String())
	})

	// QueryParam
	paramCases := []struct {
		name   string
		params string
	}{
		{
			name:   "validate error handling for query first",
			params: "first=1000",
		},
		{
			name:   "validate error handling for query order",
			params: "first=1&order=AAA",
		},
	}

	for _, tc := range paramCases {
		t.Run("Test "+tc.name, func(t *testing.T) {
			mockModel := &mockRepoModel{}
			controller := NewRepoController(mockModel)

			res := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(res)
			url := "/repo?" + tc.params
			c.Request, _ = http.NewRequest(http.MethodGet, url, nil)
			c.Request.Header.Set("Authorization", "Bearer")

			controller.IndexByToken(c)

			assert.Equal(t, http.StatusBadRequest, res.Code)
		})
	}
}
