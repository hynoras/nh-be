package app

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
)

type HealthDeps struct {
	SQLDB    *sql.DB
	Redis    *redis.Client
	RabbitMQ *amqp.Connection
}

func RegisterHealthRoutes(r *gin.Engine, deps HealthDeps) {
	health := r.Group("/health")
	{
		health.GET("/live", liveHandler())
		health.GET("/ready", readyHandler(deps))
	}
}

func liveHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "alive",
		})
	}
}

func readyHandler(deps HealthDeps) gin.HandlerFunc {
	return func(c *gin.Context) {
		ready := true
		checks := gin.H{}

		// Check Postgres
		if deps.SQLDB != nil {
			if err := deps.SQLDB.PingContext(c.Request.Context()); err != nil {
				ready = false
				checks["postgres"] = gin.H{"status": "down", "error": err.Error()}
			} else {
				checks["postgres"] = gin.H{"status": "up"}
			}
		} else {
			ready = false
			checks["postgres"] = gin.H{"status": "down", "error": "connection not initialized"}
		}

		if deps.Redis != nil {
			if err := deps.Redis.Ping(c.Request.Context()).Err(); err != nil {
				ready = false
				checks["redis"] = gin.H{"status": "down", "error": err.Error()}
			} else {
				checks["redis"] = gin.H{"status": "up"}
			}
		} else {
			ready = false
			checks["redis"] = gin.H{"status": "down", "error": "connection not initialized"}
		}

		// Check RabbitMQ
		if deps.RabbitMQ != nil {
			if deps.RabbitMQ.IsClosed() {
				ready = false
				checks["rabbitmq"] = gin.H{"status": "down", "error": "connection closed"}
			} else {
				checks["rabbitmq"] = gin.H{"status": "up"}
			}
		} else {
			ready = false
			checks["rabbitmq"] = gin.H{"status": "down", "error": "connection not initialized"}
		}

		status := http.StatusOK
		statusText := "ready"
		if !ready {
			status = http.StatusServiceUnavailable
			statusText = "not_ready"
		}

		c.JSON(status, gin.H{
			"status": statusText,
			"checks": checks,
		})
	}
}
