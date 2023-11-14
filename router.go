package main

import "net/http"

import (
	"ETM/models"
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
	apiV1 := httpRouter.Group("/api/v1")
	{
		apiV1.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "pong",
			})
		})
	}
	// Serve frontend static files

	apiV1.GET("/categories", models.GetCategories)
	apiV1.PUT("/categories", models.CreateCategory)
	apiV1.GET("/tasks/:categoryId", models.GetTasks)
	apiV1.GET("/task/:id", models.GetTask)
	apiV1.DELETE("/task/:id", models.DeleteTask)
	apiV1.POST("/task", models.CreateTask)
	apiV1.PUT("/task/:id", models.UpdateTask)
	apiV1.GET("/healthcheck", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "OK",
		})
	})

	// Start and run the server
	err := httpRouter.Run(listenAddr)
	return err
}
