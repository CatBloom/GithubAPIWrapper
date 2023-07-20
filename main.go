package main

import (
	"net/http"
	"os"
	"time"

	"main/configs"
	"main/controllers"
	"main/middleware"
	"main/models"
	"main/utils"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var r *gin.Engine

func init() {
	host := ""
	if os.Getenv("ENV") == "local" {
		utils.InitEnv()
		host = configs.GetLocalHost()
	} else {
		host = configs.GetHost()
	}

	authModel := models.NewAuthModel()
	authController := controllers.NewAuthController(authModel)

	userModel := models.NewUserModel()
	userController := controllers.NewUserController(userModel)

	repoModel := models.NewRepoModel()
	repoController := controllers.NewRepoController(repoModel)

	issueModel := models.NewIssueModel()
	issueController := controllers.NewIssueController(issueModel)

	r = gin.Default()

	// cors設定
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			host,
		},
		AllowMethods: []string{
			"GET",
			"POST",
			"PUT",
			"DELETE",
			"OPTIONS",
		},
		AllowHeaders: []string{
			"Accept", "Authorization", "Content-Type",
		},
		AllowCredentials: true,
		MaxAge:           24 * time.Hour,
	}))

	// router設定
	api := r.Group("/api")
	api.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, "Welcome Github API Wrapper!")
	})

	auth := api.Group("/auth")
	{
		auth.GET("", authController.Get)
		// local開発用のTokenを返す
		auth.GET("/test", func(c *gin.Context) {
			cookie := &http.Cookie{
				Name:     "Token",
				Value:    os.Getenv("ACCESS_TOKEN"),
				Path:     "/",
				Expires:  time.Now().Add(24 * time.Hour),
				HttpOnly: true,
			}
			http.SetCookie(c.Writer, cookie)
			c.JSON(http.StatusOK, "Github API Local Token!")
		})
	}

	viewer := api.Group("/viewer")
	viewer.Use(middleware.ParseAuthHandler())
	{
		viewer.GET("/user", userController.GetByToken)
		viewer.GET("/repos", repoController.IndexByToken)
	}

	issue := api.Group("/issue")
	issue.Use(middleware.ParseAuthHandler())
	{
		issue.GET("/list", issueController.Index)
		issue.GET("", issueController.Get)
		issue.POST("", issueController.Create)
	}
}

func main() {
	r.Run()
}
