package controllers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"
	"workship-disaster-api/models"
	"workship-disaster-api/resp"
	"workship-disaster-api/service"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

// AssignmentController ...
type AssignmentController struct {
	db                *sql.DB
	rdb               *redis.Client
	assignmentService *service.AssignmentService
}

// NewAssignmentController ...
func NewAssignmentController(db *sql.DB, rdb *redis.Client) *AssignmentController {
	areaService := service.NewAreaService(db)
	truckService := service.NewTruckService(db)
	assignmentService := service.NewAssignmentService(areaService, truckService)

	return &AssignmentController{
		db:                db,
		rdb:               rdb,
		assignmentService: assignmentService,
	}
}

// CreateAssignment handles the creation of truck assignments for areas
func (c *AssignmentController) CreateAssignment(ctx *gin.Context) {
	// Try to get from cache first
	cacheKey := "assignments:latest"
	cachedResult, err := c.rdb.Get(ctx, cacheKey).Result()
	if err == nil {
		var assignments []models.Assignment
		if err := json.Unmarshal([]byte(cachedResult), &assignments); err == nil {
			ctx.JSON(http.StatusOK, resp.SuccessResponse{
				Code:    http.StatusOK,
				Message: "Assignments retrieved from cache",
				Data:    assignments,
			})
			return
		}
	}

	// If not in cache or error, create new assignments
	assignments, err := c.assignmentService.CreateAssignments()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, resp.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to create assignments",
			Error:   err.Error(),
		})
		return
	}

	// Cache the new assignments expire time 30 mins
	if jsonData, err := json.Marshal(assignments); err == nil {
		c.rdb.Set(ctx, cacheKey, jsonData, 30*time.Minute)
	}

	ctx.JSON(http.StatusOK, resp.SuccessResponse{
		Code:    http.StatusOK,
		Message: "Assignments created successfully",
		Data:    assignments,
	})
}

// GetAssignments retrieves the latest assignments from cache
func (c *AssignmentController) GetAssignments(ctx *gin.Context) {
	cacheKey := "assignments:latest"
	cachedResult, err := c.rdb.Get(ctx, cacheKey).Result()
	if err != nil {
		ctx.JSON(http.StatusNotFound, resp.ErrorResponse{
			Code:    http.StatusNotFound,
			Message: "No assignments found in cache",
			Error:   "Cache miss",
		})
		return
	}

	var assignments []models.Assignment
	if err := json.Unmarshal([]byte(cachedResult), &assignments); err != nil {
		ctx.JSON(http.StatusInternalServerError, resp.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to parse cached assignments",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, resp.SuccessResponse{
		Code:    http.StatusOK,
		Message: "Assignments retrieved from cache",
		Data:    assignments,
	})
}

// DeleteAssignments clears the assignments from cache
func (c *AssignmentController) DeleteAssignments(ctx *gin.Context) {
	cacheKey := "assignments:latest"
	err := c.rdb.Del(ctx, cacheKey).Err()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, resp.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to clear assignments cache",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, resp.SuccessResponse{
		Code:    http.StatusOK,
		Message: "Assignments cache cleared successfully",
		Data:    nil,
	})
}
