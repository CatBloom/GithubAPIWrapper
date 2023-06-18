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

	userModeles := models.NewUserModels()
	userController := controllers.NewUserController(userModeles)

	repoModeles := models.NewRepoModels()
	repoController := controllers.NewRepoController(repoModeles)

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

	user := api.Group("/user")
	{
		user.GET("", userController.GetByToken)
	}

	repo := api.Group("/repo")
	{
		repo.GET("", repoController.IndexByToken)
	}
}

func main() {
	r.Run()
}
