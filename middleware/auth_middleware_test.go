package middleware

import (
	"main/utils"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestParseAuthHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	utils.InitEnv()
	t.Run("Test with valid token", func(t *testing.T) {
		res := httptest.NewRecorder()
		c, router := gin.CreateTestContext(res)

		cookie := &http.Cookie{
			Name:  "Token",
			Value: os.Getenv("ACCESS_TOKEN"), Path: "/",
			Expires: time.Now().Add(24 * time.Hour),
		}
		c.Request, _ = http.NewRequest("GET", "/", nil)
		c.Request.AddCookie(cookie)

		router.Use(ParseAuthHandler())
		router.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, "")
		})

		router.ServeHTTP(res, c.Request)

		assert.Equal(t, http.StatusOK, res.Code)
	})

	t.Run("Test with token is missing in cookie", func(t *testing.T) {
		res := httptest.NewRecorder()
		c, router := gin.CreateTestContext(res)

		cookie := &http.Cookie{}
		c.Request, _ = http.NewRequest("GET", "/", nil)
		c.Request.AddCookie(cookie)

		router.Use(ParseAuthHandler())
		router.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, "")
		})

		router.ServeHTTP(res, c.Request)

		assert.Equal(t, http.StatusAccepted, res.Code)
		expectedResponse := `{"message":"token is missing in cookie"}`
		assert.Equal(t, expectedResponse, res.Body.String())
	})

	t.Run("Test with token is empty in cookie", func(t *testing.T) {
		res := httptest.NewRecorder()
		c, router := gin.CreateTestContext(res)

		cookie := &http.Cookie{
			Name:    "Token",
			Value:   "",
			Expires: time.Now().Add(24 * time.Hour),
		}
		c.Request, _ = http.NewRequest("GET", "/", nil)
		c.Request.AddCookie(cookie)

		router.Use(ParseAuthHandler())
		router.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, "")
		})

		router.ServeHTTP(res, c.Request)

		assert.Equal(t, http.StatusAccepted, res.Code)
		expectedResponse := `{"message":"token is empty in cookie"}`
		assert.Equal(t, expectedResponse, res.Body.String())
	})
}
