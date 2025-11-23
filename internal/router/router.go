package router

import (
	"github.com/askaroe/dockify-backend/internal/handlers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewRouter(handler *handlers.Handler) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Replace with specific origins if needed
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"}, // Replace with specific headers if needed
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum cache age in seconds
	}))

	api := r.Group("/api/v1")
	{
		api.POST("/register", handler.Register)
		api.POST("/login", handler.Login)
		api.POST("/metrics", handler.Health.CreateHealthMetrics)
		api.GET("/metrics", handler.Health.GetHealthMetrics)
		api.GET("/recommendation", handlers.GetRecommendation)

		location := api.Group("/location")
		{
			location.POST("/nearest", handler.Location.GetNearestUsers)
		}

		hospitals := api.Group("/hospitals")
		{
			hospitals.POST("/nearest", handlers.GetNearestHospitals)
		}
	}

	r.GET("/health", handlers.HealthCheck)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}

func SetJSONContentType() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Next()
	}
}
