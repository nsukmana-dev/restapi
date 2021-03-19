package main

import (
	"github.com/gin-gonic/gin"
	"github.com/nsukmana-dev/restapi/config"
	"github.com/nsukmana-dev/restapi/routes"
	"github.com/subosito/gotenv"
)

func main() {
	config.InitDB()
	defer config.DB.Close()
	gotenv.Load()

	router := gin.Default()

	v1 := router.Group("/api/v1/")
	{
		v1.GET("/auth/:provider/", routes.RedirectHandler)
		v1.GET("/auth/:provider/callback/", routes.CallbackHandler)

		articles := v1.Group("/article")
		{
			articles.GET("/", routes.GetHome)
			articles.GET("/:slug", routes.GetArticle)
			articles.POST("/", routes.PostArticle)
		}
	}

	router.Run()
}