package webserver

import (
	"ETM/pkg/app"
	controllers2 "ETM/pkg/controllers"
	"ETM/pkg/vars"
	"net/http"

	"github.com/gin-gonic/contrib/cors"
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
	if vars.AllowedOrigin == "*" {
		config.AllowAllOrigins = true
		config.AllowedOrigins = nil
	} else {
		config.AllowedOrigins = []string{vars.AllowedOrigin}
		config.AllowAllOrigins = false
	}
	config.AddAllowedMethods("OPTIONS")
	config.AllowedHeaders = []string{"Authorization", "Content-Type"}
	httpRouter.Use(cors.New(config))

	httpRouter.Use(AppHandler(App))

	// Serve compiled frontend assets directly; everything else falls to NoRoute.
	httpRouter.Static("/assets", "./yaftm/dist/assets")
	httpRouter.StaticFile("/favicon.ico", "./yaftm/dist/favicon.ico")

	// SPA fallback: any unmatched path (including /, /login, /auth/callback …)
	// gets index.html so Vue Router can take over client-side.
	httpRouter.NoRoute(func(c *gin.Context) {
		c.File("./yaftm/dist/index.html")
	})
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
	apiV1.GET("/categories", controllers2.IsAuthorized(), controllers2.GetCategories)
	apiV1.POST("/category", controllers2.IsAuthorized(), controllers2.CreateCategory)
	apiV1.POST("/category/:categoryId/share", controllers2.IsAuthorized(), controllers2.ShareCategory)
	apiV1.DELETE("/category/:categoryId/share/:groupId", controllers2.IsAuthorized(), controllers2.UnshareCategory)
	apiV1.GET("/category/:categoryId/shares", controllers2.IsAuthorized(), controllers2.GetCategoryShares)

	// Group endpoints
	apiV1.POST("/group", controllers2.IsAuthorized(), controllers2.CreateGroup)
	apiV1.GET("/groups", controllers2.IsAuthorized(), controllers2.GetGroups)
	apiV1.GET("/group/:groupId", controllers2.IsAuthorized(), controllers2.GetGroup)
	apiV1.DELETE("/group/:groupId", controllers2.IsAuthorized(), controllers2.DeleteGroup)
	apiV1.POST("/group/:groupId/member", controllers2.IsAuthorized(), controllers2.AddGroupMember)
	apiV1.PUT("/group/:groupId/member/:userId", controllers2.IsAuthorized(), controllers2.UpdateGroupMember)
	apiV1.DELETE("/group/:groupId/member/:userId", controllers2.IsAuthorized(), controllers2.RemoveGroupMember)

	// Tasks Endpoints

	apiV1.GET("/tasks/:categoryId", controllers2.IsAuthorized(), controllers2.GetTasks)
	apiV1.POST("/task", controllers2.IsAuthorized(), controllers2.CreateTask)
	apiV1.POST("/task/:taskId", controllers2.IsAuthorized(), controllers2.UpdateTask)
	apiV1.GET("/task/:taskId", controllers2.IsAuthorized(), controllers2.GetTask)
	apiV1.DELETE("/task/:taskId", controllers2.IsAuthorized(), controllers2.DeleteTask)

	// User endpoints
	apiV1.GET("/user/logout", controllers2.Logout)
	apiV1.POST("/user/login", controllers2.Login)
	apiV1.POST("/user/register", controllers2.Register)
	apiV1.POST("/user", controllers2.IsAuthorized(), controllers2.UpdateUser)
	apiV1.GET("/user", controllers2.IsAuthorized(), controllers2.GetUser)
	apiV1.GET("/user/device", controllers2.IsAuthorized(), controllers2.GetUserDevices)
	apiV1.GET("/user/refreshtoken", controllers2.IsAuthorized(), controllers2.RefreshToken)
	apiV1.POST("/user/updatesubscription", controllers2.IsAuthorized(), controllers2.UpdateUserSubscription)
	apiV1.GET("/user/devices", controllers2.IsAuthorized(), controllers2.GetUserDevices)
	apiV1.GET("/user/tokens", controllers2.IsAuthorized(), controllers2.GetUserTokens)
	apiV1.POST("/user/devices", controllers2.IsAuthorized(), controllers2.CreateUserDevice)

	// OIDC endpoints
	apiV1.GET("/auth/oidc/status", controllers2.OIDCStatus)
	apiV1.GET("/auth/oidc/login", controllers2.OIDCLogin)
	apiV1.GET("/auth/oidc/callback", controllers2.OIDCCallback)

	apiV1.GET("/healthcheck", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "OK",
		})
	})

	// Start and run the server
	err := httpRouter.Run(listenAddr)
	return err
}
