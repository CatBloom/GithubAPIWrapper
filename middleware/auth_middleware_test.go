package middleware

import (
	"main/utils"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestParseAuthHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	utils.InitEnv()
	t.Run("Test with valid token", func(t *testing.T) {
		res := httptest.NewRecorder()
		c, router := gin.CreateTestContext(res)

		router.Use(ParseAuthHandler())
		router.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, "")
		})
		c.Request, _ = http.NewRequest("GET", "/", nil)
		testJwtStr := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbiI6InRva2VuLXNhbXBsZSJ9.GNM7Iw3xiwEcvXXN9RYM9fFFGl06dKb7Q0sCKN3--lw"
		c.Request.Header.Set("Authorization", testJwtStr)

		router.ServeHTTP(res, c.Request)

		assert.Equal(t, http.StatusOK, res.Code)
	})

	t.Run("Test with missing token", func(t *testing.T) {
		res := httptest.NewRecorder()
		c, router := gin.CreateTestContext(res)

		router.Use(ParseAuthHandler())
		router.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, "")
		})
		c.Request, _ = http.NewRequest("GET", "/", nil)

		router.ServeHTTP(res, c.Request)

		assert.Equal(t, http.StatusUnauthorized, res.Code)
		expectedResponse := `{"error":"Error invalid authorization token"}`
		assert.Equal(t, expectedResponse, res.Body.String())
	})

	t.Run("Test with invalid token", func(t *testing.T) {
		res := httptest.NewRecorder()
		c, router := gin.CreateTestContext(res)

		router.Use(ParseAuthHandler())
		router.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, "")
		})
		c.Request, _ = http.NewRequest("GET", "/", nil)
		c.Request.Header.Set("Authorization", "123456")

		router.ServeHTTP(res, c.Request)

		assert.Equal(t, http.StatusUnauthorized, res.Code)
		expectedResponse := `{"error":"token contains an invalid number of segments"}`
		assert.Equal(t, expectedResponse, res.Body.String())
	})
}
