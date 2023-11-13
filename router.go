package main

import "net/http"

import (
	"ETM/models"
	"github.com/gin-gonic/gin"
)

func runHttp(listenAddr string) error {

	httpRouter := gin.Default()

	apiV1 := httpRouter.Group("/api/v1")
	{
		apiV1.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "pong",
			})
		})
	}
	apiV1.GET("/categories", models.GetCategories)
	apiV1.PUT("/categories", models.CreateCategory)
	apiV1.GET("/tasks", models.GetTasks)
	apiV1.GET("/tasks/:id", models.GetTask)
	apiV1.DELETE("/tasks/:id", models.DeleteTask)
	apiV1.POST("/tasks", models.CreateTask)
	apiV1.PUT("/tasks/:id", models.UpdateTask)
	apiV1.GET("/healthcheck", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "OK",
		})
	})

	// Start and run the server
	err := httpRouter.Run(listenAddr)
	return err
}
