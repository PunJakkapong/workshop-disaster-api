package controllers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"
	"workship-disaster-api/models"
	"workship-disaster-api/resp"

	"github.com/gin-gonic/gin"
)

type AreaController struct {
	db *sql.DB
}

func NewAreaController(db *sql.DB) *AreaController {
	return &AreaController{db: db}
}

// CreateArea handles the creation of a new area
func (c *AreaController) CreateArea(ctx *gin.Context) {
	var req models.CreateAreaRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, resp.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Binding error.",
			Error:   err.Error(),
		})
		return
	}

	// Check if area ID already exists
	var exists bool
	err := c.db.QueryRow("SELECT EXISTS(SELECT 1 FROM areas WHERE area_id = $1)", req.AreaID).Scan(&exists)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, resp.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to check area ID",
			Error:   err.Error(),
		})
		return
	}
	if exists {
		ctx.JSON(http.StatusConflict, resp.ErrorResponse{
			Code:    http.StatusConflict,
			Message: "Area ID already exists",
		})
		return
	}

	// Convert RequiredResources to JSON
	resourcesJSON, err := json.Marshal(req.RequiredResources)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, resp.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to process required resources",
			Error:   err.Error(),
		})
		return
	}

	// Insert into database
	_, err = c.db.Exec(
		"INSERT INTO areas (area_id, urgency_level, required_resources, time_constraint) VALUES ($1, $2, $3, $4)",
		req.AreaID, req.UrgencyLevel, resourcesJSON, req.TimeConstraint,
	)

	if err != nil {
		// Check if error is due to duplicate key
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			ctx.JSON(http.StatusConflict, resp.ErrorResponse{
				Code:    http.StatusConflict,
				Message: "Area ID already exists",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, resp.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to create area",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, resp.SuccessResponse{
		Code:    http.StatusCreated,
		Message: "Area created successfully",
		Data: gin.H{
			"areaId": req.AreaID,
		},
	})
}
