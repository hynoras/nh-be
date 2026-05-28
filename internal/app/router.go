package app

import (
	"net/http"

	docs "nh-be/docs"
	"nh-be/internal/config"
	"nh-be/internal/middleware"

	"nh-be/internal/constant"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

// getAllowedOrigins returns a slice of allowed CORS origins based on the application
// environment. In development mode ("dev"), only localhost:3000 is allowed.
// In production or other environments, both localhost:3000 and the configured
// FrontendURL are permitted.
func getAllowedOrigins(cfg *config.Config) []string {
	if cfg.AppEnv == "dev" {
		return []string{"http://localhost:3000"}
	}
	return []string{"http://localhost:3000", cfg.FrontendURL}
}

// getCorsConfig constructs the CORS middleware configuration using the provided
// allowed origins and the predefined security constants from the constant package.
func getCorsConfig(origins []string) cors.Config {
	return cors.Config{
		AllowOrigins:     origins,
		AllowMethods:     constant.CorsAllowMethods,
		AllowHeaders:     constant.CorsAllowHeaders,
		ExposeHeaders:    constant.CorsExposeHeaders,
		AllowCredentials: constant.CorsAllowCredentials,
		MaxAge:           constant.CorsMaxAge,
	}
}

// NewRouter creates and configures the Gin engine with the full middleware
// chain (Recovery, CORS, OTEL tracing, request ID, structured logging, and
// Prometheus metrics), registers Swagger documentation, and utility routes
// (/ping, /metrics).
//
// Feature routes are NOT registered here to avoid an import cycle between
// internal/app and internal/router. The caller is responsible for calling
// router.SetupRoutes after obtaining the engine.
func NewRouter(cfg *config.Config) *gin.Engine {
	allowedOrigins := getAllowedOrigins(cfg)
	corsConfig := getCorsConfig(allowedOrigins)

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(cors.New(corsConfig))

	r.Use(otelgin.Middleware("noheir-api"))
	r.Use(middleware.SetRequestID())
	r.Use(middleware.RequestLogger())
	r.Use(middleware.MetricsMiddleware())

	docs.SwaggerInfo.BasePath = "/api/v1"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	return r
}
