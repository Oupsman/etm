package webserver

import (
	"ETM/pkg/app"
	controllers2 "ETM/pkg/controllers"
	"ETM/pkg/models"
	"github.com/gin-gonic/contrib/cors"
	"net/http"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
)

func AppHandler(App *app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("App", App)
		c.Next()
	}
}

func RunHttp(listenAddr string, App *app.App) error {

	httpRouter := gin.Default()

	//	httpRouter.LoadHTMLGlob("templates/*")
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AddAllowedMethods("OPTIONS")
	config.AllowedHeaders = []string{"Authorization", "Content-Type"}
	httpRouter.Use(cors.New(config))

	httpRouter.Use(static.Serve("/static", static.LocalFile("./static", true)))
	httpRouter.Use(AppHandler(App))
	httpRouter.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title": "Eisenhower TaskCard Manager",
		})
	})

	httpRouter.GET("/signup", func(c *gin.Context) {
		c.HTML(http.StatusOK, "signup.tmpl", gin.H{
			"title": "Eisenhower TaskCard Manager - Sign up",
		})
	})
	// Serve frontend static files

	httpRouter.StaticFile("/service-worker.js", "./resources/service-worker.js")

	apiV1 := httpRouter.Group("/api/v1")
	{
		apiV1.GET("/", controllers2.IsAuthorized(), func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "pong",
			})
		})
	}

	// Generic endpoints

	apiV1.GET("/getvapidkey", controllers2.IsAuthorized(), controllers2.GetVAPIDKey)

	// Categories endpoints
	apiV1.GET("/categories", models.GetCategories)
	apiV1.POST("/categories", models.CreateCategory)

	// Tasks Endpoints

	apiV1.GET("/tasks/:categorieId", controllers2.GetTasks)
	apiV1.POST("/task", controllers2.CreateTask)
	apiV1.POST("/task/:taskId", controllers2.UpdateTask)
	apiV1.GET("/task/:taskId", controllers2.GetTask)
	apiV1.DELETE("/task/:taskId", controllers2.DeleteTask)

	// User endpoints
	apiV1.GET("/user/logout", controllers2.Logout)
	apiV1.POST("/user/login", controllers2.Login)
	apiV1.POST("/user/register", controllers2.Register)
	apiV1.POST("/user", controllers2.UpdateUser)
	apiV1.GET("/user", controllers2.GetUser)
	apiV1.GET("/user/refreshtoken", controllers2.RefreshToken)
	apiV1.POST("/user/updatesubscription", controllers2.UpdateUserSubscription)

	apiV1.GET("/healthcheck", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "OK",
		})
	})

	// Start and run the server
	err := httpRouter.Run(listenAddr)
	return err
}
