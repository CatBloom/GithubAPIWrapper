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

type mockIssueModel struct{}

func (m *mockIssueModel) GetIssues(_ string, _ types.IssuesReq) (types.IssuesResp, error) {
	issues := []types.Issue{
		{
			ID:        "IssueID1",
			CreatedAt: "2000-01-01T00:00:00Z",
			UpdatedAt: "2000-01-01T00:00:00Z",
			State:     "OPEN",
			URL:       "https://example.com/issue",
			Title:     "First issue",
			Number:    1,
			Body:      "Body issue",
			BodyHTML:  "Body issue",
		},
		{
			ID:        "IssueID2",
			CreatedAt: "2000-01-01T00:00:00Z",
			UpdatedAt: "2000-01-01T00:00:00Z",
			State:     "OPEN",
			URL:       "https://example.com/issue",
			Title:     "Second issue",
			Number:    2,
			Body:      "Body issue",
			BodyHTML:  "Body issue",
		},
	}

	resp := types.IssuesResp{
		Data: types.IssuesRepository{
			Repository: types.Issues{
				CreatedAt:     "2000-01-01T00:00:00Z",
				UpdatedAt:     "2000-01-01T00:00:00Z",
				Name:          "RepoName",
				NameWithOwner: "Owner/RepoName",
				Issues: types.IssuesNodes{
					Nodes: issues,
				},
			},
		},
	}

	return resp, nil
}

func (m *mockIssueModel) GetIssue(_ string, _ types.IssueReq) (types.IssueResp, error) {
	issue := types.Issue{
		ID:        "IssueID1",
		CreatedAt: "2000-01-01T00:00:00Z",
		UpdatedAt: "2000-01-01T00:00:00Z",
		State:     "OPEN",
		URL:       "https://example.com/issue",
		Title:     "First issue",
		Number:    1,
		Body:      "Body issue",
		BodyHTML:  "Body issue",
	}

	resp := types.IssueResp{
		Data: types.IssueRepository{
			Repository: types.IssueNode{
				Issue: issue,
			},
		},
	}

	return resp, nil
}

type mockErrorIssueModel struct{}

func (m *mockErrorIssueModel) GetIssues(_ string, _ types.IssuesReq) (types.IssuesResp, error) {
	errorMessage := "Failed to issues"
	return types.IssuesResp{}, errors.New(errorMessage)
}

func (m *mockErrorIssueModel) GetIssue(_ string, _ types.IssueReq) (types.IssueResp, error) {
	errorMessage := "Failed to issue"
	return types.IssueResp{}, errors.New(errorMessage)
}

func TestIndexIssueController(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Run("Test with success response", func(t *testing.T) {
		mockModel := &mockIssueModel{}
		controller := NewIssueController(mockModel)

		res := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(res)
		c.Request, _ = http.NewRequest(http.MethodGet, "/issue/list?first=1&order=DESC&owner=aaa&repo=bbb&states=OPEN", nil)
		c.Set("token", "token")
		c.Request.Header.Set("Authorization", "Bearer")

		controller.Index(c)

		assert.Equal(t, http.StatusOK, res.Code)
	})

	t.Run("Test with error response", func(t *testing.T) {
		mockErrorModel := &mockErrorIssueModel{}
		controller := NewIssueController(mockErrorModel)

		res := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(res)
		c.Request, _ = http.NewRequest(http.MethodGet, "/issue/list?first=1&order=DESC&owner=aaa&repo=bbb&states=OPEN", nil)
		c.Request.Header.Set("Authorization", "Bearer")

		controller.Index(c)

		assert.Equal(t, http.StatusNotFound, res.Code)
		expectedResponse := `{"error":"Failed to issues"}`
		assert.Equal(t, expectedResponse, res.Body.String())
	})

	// QueryParam
	paramCases := []struct {
		name   string
		params string
	}{
		{
			name:   "validate error handling for query first",
			params: "first=0&order=DESC&owner=aaa&repo=bbb&states=OPEN",
		},
		{
			name:   "validate error handling for query order",
			params: "first=1&order=AAAA&owner=aaa&repo=bbb&states=OPEN",
		},
		{
			name:   "validate error handling for query owner",
			params: "first=1&order=DESC&repo=bbb&states=OPEN",
		},
		{
			name:   "validate error handling for query repo",
			params: "first=1&order=DESC&owner=aaa&states=OPEN",
		},
		{
			name:   "validate error handling for query states",
			params: "first=1&order=DESC&owner=aaa&repo=bbb&states=AAA",
		},
	}

	for _, tc := range paramCases {
		t.Run("Test "+tc.name, func(t *testing.T) {
			mockModel := &mockIssueModel{}
			controller := NewIssueController(mockModel)

			res := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(res)
			url := "/issue/list?" + tc.params
			c.Request, _ = http.NewRequest(http.MethodGet, url, nil)
			c.Request.Header.Set("Authorization", "Bearer")

			controller.Index(c)

			assert.Equal(t, http.StatusBadRequest, res.Code)
		})
	}

}

func TestGetIssueController(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Run("Test with success response", func(t *testing.T) {
		mockModel := &mockIssueModel{}
		controller := NewIssueController(mockModel)

		res := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(res)
		c.Request, _ = http.NewRequest(http.MethodGet, "/issue?&owner=aaa&repo=bbb&number=1", nil)
		c.Set("token", "token")
		c.Request.Header.Set("Authorization", "Bearer")

		controller.Get(c)

		assert.Equal(t, http.StatusOK, res.Code)
	})

	t.Run("Test with error response", func(t *testing.T) {
		mockErrorModel := &mockErrorIssueModel{}
		controller := NewIssueController(mockErrorModel)

		res := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(res)
		c.Request, _ = http.NewRequest(http.MethodGet, "/issue?&owner=aaa&repo=bbb&number=1", nil)
		c.Request.Header.Set("Authorization", "Bearer")

		controller.Get(c)

		assert.Equal(t, http.StatusNotFound, res.Code)
		expectedResponse := `{"error":"Failed to issue"}`
		assert.Equal(t, expectedResponse, res.Body.String())
	})

	// QueryParam
	paramCases := []struct {
		name   string
		params string
	}{
		{
			name:   "validate error handling for query order",
			params: "owner=aaa&repo=bbb&number=",
		},
		{
			name:   "validate error handling for query owner",
			params: "owner=&repo=bbb&number=1",
		},
		{
			name:   "validate error handling for query repo",
			params: "owner=aaa&repo=&number=1",
		},
	}

	for _, tc := range paramCases {
		t.Run("Test "+tc.name, func(t *testing.T) {
			mockModel := &mockIssueModel{}
			controller := NewIssueController(mockModel)

			res := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(res)
			url := "/issue/list?" + tc.params
			c.Request, _ = http.NewRequest(http.MethodGet, url, nil)
			c.Request.Header.Set("Authorization", "Bearer")

			controller.Get(c)

			assert.Equal(t, http.StatusBadRequest, res.Code)
		})
	}

}
