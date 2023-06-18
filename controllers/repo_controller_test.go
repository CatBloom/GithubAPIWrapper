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

type mockRepoModels struct{}

func (m *mockRepoModels) GetReposByToken(_ string, first int, order string) (models.RepoResp, error) {
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
				},
			},
		},
	}

	repoResp := resp
	return repoResp, nil
}

type mockErrorRepoModels struct{}

func (m *mockErrorRepoModels) GetReposByToken(_ string, first int, order string) (models.RepoResp, error) {
	errorMessage := "Failed to repositories"
	return models.RepoResp{}, errors.New(errorMessage)
}

func TestIndexByTokenControllers(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Run("Test with success response", func(t *testing.T) {
		mockModels := &mockRepoModels{}
		controller := NewRepoController(mockModels)

		res := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(res)
		c.Request, _ = http.NewRequest(http.MethodGet, "/repo?first=1&order=DESC", nil)
		c.Set("token", "token")
		c.Request.Header.Set("Authorization", "Bearer")

		controller.IndexByToken(c)

		assert.Equal(t, http.StatusOK, res.Code)
		expectedResponse := `{"data":{"viewer":{"repositories":{"nodes":[{"name":"First Repository","description":"description","url":"https://example.com/repository","createdAt":"2000-01-01T00:00:00Z","updatedAt":"2000-01-01T00:00:00Z"},{"name":"Second Repository","description":"","url":"https://example.com/repository","createdAt":"2000-01-01T00:00:00Z","updatedAt":"2000-01-01T00:00:00Z"}]}}}}`
		assert.Equal(t, expectedResponse, res.Body.String())
	})

	t.Run("Test validate error handling for query first", func(t *testing.T) {
		mockModels := &mockRepoModels{}
		controller := NewRepoController(mockModels)

		res := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(res)
		c.Request, _ = http.NewRequest(http.MethodGet, "/repo?first=1000&order=DESC", nil)
		c.Request.Header.Set("Authorization", "Bearer")

		controller.IndexByToken(c)

		assert.Equal(t, http.StatusBadRequest, res.Code)
	})

	t.Run("Test validate error handling for query order", func(t *testing.T) {
		mockModels := &mockRepoModels{}
		controller := NewRepoController(mockModels)

		res := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(res)
		c.Request, _ = http.NewRequest(http.MethodGet, "/repo?first=1&order=AAA", nil)
		c.Request.Header.Set("Authorization", "Bearer")

		controller.IndexByToken(c)

		assert.Equal(t, http.StatusBadRequest, res.Code)
	})

	t.Run("Test with error response", func(t *testing.T) {
		mockErrorModels := &mockErrorRepoModels{}
		controller := NewRepoController(mockErrorModels)

		res := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(res)
		c.Request, _ = http.NewRequest(http.MethodGet, "/repo?first=1&order=DESC", nil)
		c.Request.Header.Set("Authorization", "Bearer")

		controller.IndexByToken(c)

		assert.Equal(t, http.StatusNotFound, res.Code)
		expectedResponse := `{"error":"Failed to repositories"}`
		assert.Equal(t, expectedResponse, res.Body.String())
	})
}
