package main

import (
	"net/http"
	"os"
	"time"

	"main/controllers"
	"main/models"
	"main/utils"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var r *gin.Engine

func init() {
	if os.Getenv("ENV") == "dev" {
		utils.InitEnv()
	}

	userModel := models.NewUserModel()
	userController := controllers.NewUserController(userModel)

	repoModel := models.NewRepoModel()
	repoController := controllers.NewRepoController(repoModel)

	r = gin.Default()

	// cors設定
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"*",
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

	viewer := api.Group("/viewer")
	{
		viewer.GET("/user", userController.GetByToken)
		viewer.GET("/repos", repoController.IndexByToken)
	}
}

func main() {
	r.Run()
}
