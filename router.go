package main

import (
	"ETM/controllers"
	"ETM/models"
	"net/http"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
)

func runHttp(listenAddr string) error {

	httpRouter := gin.Default()

	httpRouter.LoadHTMLGlob("templates/*")

	httpRouter.Use(static.Serve("/static", static.LocalFile("./static", true)))
	httpRouter.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title": "Eisenhower Task Manager",
		})
	})

	httpRouter.GET("/signup", func(c *gin.Context) {
		c.HTML(http.StatusOK, "signup.tmpl", gin.H{
			"title": "Eisenhower Task Manager - Sign up",
		})
	})
	// Serve frontend static files

	httpRouter.StaticFile("/service-worker.js", "./resources/service-worker.js")

	apiV1 := httpRouter.Group("/api/v1")
	{
		apiV1.GET("/", controllers.IsAuthorized(), func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "pong",
			})
		})
	}

	// Categories endpoints
	apiV1.GET("/categories", models.GetCategories)
	apiV1.POST("/categories", models.CreateCategory)

	// Tasks Endpoints

	apiV1.GET("/tasks/:categoryId", controllers.GetTasks)
	apiV1.POST("/task", controllers.CreateTask)
	apiV1.POST("/task/:taskId", controllers.UpdateTask)
	apiV1.GET("/task/:taskId", controllers.GetTask)
	apiV1.GET("/task/:taskId/delete", controllers.DeleteTask)

	// User endpoints
	apiV1.GET("/user/logout", controllers.Logout)
	apiV1.POST("/user/login", controllers.Login)
	apiV1.POST("/user/register", controllers.Register)
	apiV1.POST("/user", controllers.UpdateUser)
	apiV1.GET("/user", controllers.GetUser)
	apiV1.GET("/user/refreshtoken", controllers.RefreshToken)

	apiV1.GET("/healthcheck", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "OK",
		})
	})

	// Start and run the server
	err := httpRouter.Run(listenAddr)
	return err
}
