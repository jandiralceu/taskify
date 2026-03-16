// Package routes wires HTTP handlers to Gin engine routes and applies
// global middleware such as CORS, tracing, and authentication.
package routes

import (
	"github.com/casbin/casbin/v3"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/jandiralceu/taskify/docs" // imported so swagger can read embedded docs
	"github.com/jandiralceu/taskify/internal/config"
	"github.com/jandiralceu/taskify/internal/handlers"
	"github.com/jandiralceu/taskify/internal/middleware"
	"github.com/jandiralceu/taskify/internal/pkg"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

// RouteConfig holds all handler dependencies required to register API routes.
type RouteConfig struct {
	AuthHandler *handlers.AuthHandler
	UserHandler *handlers.UserHandler
	TaskHandler *handlers.TaskHandler
}

// Setup creates a configured [gin.Engine] with global middleware, public and
// protected route groups, and the Swagger UI endpoint.
func Setup(routeConfig *RouteConfig, config *config.Config, jwtManager *pkg.JWTManager, enforcer *casbin.Enforcer, cacheManager pkg.CacheManager) *gin.Engine {
	gin.ForceConsoleColor()

	router := gin.New()
	router.Use(middleware.TraceIDMiddleware())

	// Global Rate Limit per IP/User
	router.Use(middleware.RateLimiter(cacheManager, "global", config.RateLimitGlobal))

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Configure CORS policy for cross-origin requests.
	router.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:    []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:   []string{"Content-Length", "X-Trace-ID"},
		MaxAge:          12 * 3600,
	}))

	router.SetTrustedProxies(nil)

	router.Use(otelgin.Middleware(config.AppName))

	api := router.Group("/api/v1")
	{
		// Public routes (no authentication required).
		auth := api.Group("/auth")
		// Strict limit for auth
		auth.Use(middleware.RateLimiter(cacheManager, "auth", config.RateLimitAuth))
		{
			auth.POST("/signin", routeConfig.AuthHandler.SignIn)
			auth.POST("/register", routeConfig.AuthHandler.Register)
			auth.POST("/refresh", routeConfig.AuthHandler.RefreshToken)
			auth.POST("/signout", routeConfig.AuthHandler.SignOut)
		}

		// Protected routes (authentication required).
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware(jwtManager))
		protected.Use(middleware.CasbinMiddleware(enforcer))
		{
			users := protected.Group("/users")
			{
				users.GET("/permissions", routeConfig.UserHandler.GetPermissions)
				users.GET("/profile", routeConfig.UserHandler.GetProfile)
				users.GET("", routeConfig.UserHandler.FindAllUsers)
				users.GET("/:id", routeConfig.UserHandler.FindUserByID)
				users.PATCH("/profile", routeConfig.UserHandler.UpdateUser)
				users.PATCH("/change-password", middleware.RateLimiter(cacheManager, "pass-change", config.RateLimitAuth), routeConfig.UserHandler.ChangePassword)
				users.POST("/avatar", routeConfig.UserHandler.UpdateAvatar)
				users.DELETE("/:id", routeConfig.UserHandler.DeleteUser)
			}

			tasks := protected.Group("/tasks")
			{
				tasks.POST("", routeConfig.TaskHandler.CreateTask)
				tasks.GET("", routeConfig.TaskHandler.ListTasks)
				tasks.GET("/:id", routeConfig.TaskHandler.GetTask)
				tasks.PATCH("/:id", routeConfig.TaskHandler.UpdateTask)
				tasks.DELETE("/:id", routeConfig.TaskHandler.DeleteTask)

				tasks.POST("/:id/notes", routeConfig.TaskHandler.AddNote)
				tasks.GET("/:id/notes", routeConfig.TaskHandler.GetNotes)
				tasks.PATCH("/notes/:noteId", routeConfig.TaskHandler.UpdateNote)
				tasks.DELETE("/notes/:noteId", routeConfig.TaskHandler.DeleteNote)

				tasks.POST("/:id/attachments", routeConfig.TaskHandler.AddAttachment)
				tasks.GET("/:id/attachments", routeConfig.TaskHandler.GetAttachments)
				tasks.DELETE("/attachments/:attachmentId", routeConfig.TaskHandler.DeleteAttachment)
			}
		}
	}

	// Swagger UI route.
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Static files for attachments
	router.Static("/uploads", config.UploadPath)

	return router
}
