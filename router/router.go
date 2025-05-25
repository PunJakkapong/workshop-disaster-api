package router

import (
	"context"
	"database/sql"

	"workship-disaster-api/controllers"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

// SetupRouter all the routes
func SetupRouter(db *sql.DB, rdb *redis.Client) *gin.Engine {
	r := gin.Default()

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Server is running!"})
	})

	// Redis test
	r.GET("/redis-test", func(c *gin.Context) {
		rdb.Set(ctx, "status", "Redis is working!", 0)
		status, _ := rdb.Get(ctx, "status").Result()
		c.JSON(200, gin.H{"redis": status})
	})

	// PostgreSQL test
	r.GET("/postgres-test", func(c *gin.Context) {
		var dbVersion string
		db.QueryRow("SELECT version();").Scan(&dbVersion)
		c.JSON(200, gin.H{"postgres": dbVersion})
	})

	// Initialize controllers
	areaController := controllers.NewAreaController(db)
	truckController := controllers.NewTruckController(db)
	assignmentController := controllers.NewAssignmentController(db)

	// API routes
	api := r.Group("/api")
	{
		// Areas
		areas := api.Group("/areas")
		{
			areas.POST("", areaController.CreateArea)
			areas.GET("", areaController.ListAreas)
			areas.GET("/:id", areaController.GetArea)
		}

		// Trucks
		trucks := api.Group("/trucks")
		{
			trucks.POST("", truckController.CreateTruck)
		}

		// Assignments
		api.POST("/assignments", assignmentController.CreateAssignment)
	}

	return r
}
