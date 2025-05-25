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

// GetArea retrieves an area by ID
func (c *AreaController) GetArea(ctx *gin.Context) {
	areaID := ctx.Param("id")

	var area models.Area
	var resourcesJSON []byte
	err := c.db.QueryRow(
		"SELECT area_id, urgency_level, required_resources, time_constraint FROM areas WHERE area_id = $1",
		areaID,
	).Scan(&area.AreaID, &area.UrgencyLevel, &resourcesJSON, &area.TimeConstraint)

	if err == sql.ErrNoRows {
		ctx.JSON(http.StatusNotFound, resp.ErrorResponse{
			Code:    http.StatusNotFound,
			Message: "Area not found",
		})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, resp.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to retrieve area",
			Error:   err.Error(),
		})
		return
	}

	// Parse RequiredResources JSON
	if err := json.Unmarshal(resourcesJSON, &area.RequiredResources); err != nil {
		ctx.JSON(http.StatusInternalServerError, resp.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to parse required resources",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, resp.SuccessResponse{
		Code:    http.StatusOK,
		Message: "Area retrieved successfully",
		Data:    area,
	})
}

// ListAreas retrieves all areas
func (c *AreaController) ListAreas(ctx *gin.Context) {
	rows, err := c.db.Query("SELECT area_id, urgency_level, required_resources, time_constraint FROM areas")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, resp.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to retrieve areas",
			Error:   err.Error(),
		})
		return
	}
	defer rows.Close()

	var areas []models.Area
	for rows.Next() {
		var area models.Area
		var resourcesJSON []byte
		if err := rows.Scan(&area.AreaID, &area.UrgencyLevel, &resourcesJSON, &area.TimeConstraint); err != nil {
			ctx.JSON(http.StatusInternalServerError, resp.ErrorResponse{
				Code:    http.StatusInternalServerError,
				Message: "Failed to scan area",
				Error:   err.Error(),
			})
			return
		}

		// Parse RequiredResources JSON
		if err := json.Unmarshal(resourcesJSON, &area.RequiredResources); err != nil {
			ctx.JSON(http.StatusInternalServerError, resp.ErrorResponse{
				Code:    http.StatusInternalServerError,
				Message: "Failed to parse required resources",
				Error:   err.Error(),
			})
			return
		}

		areas = append(areas, area)
	}

	ctx.JSON(http.StatusOK, resp.SuccessResponse{
		Code:    http.StatusOK,
		Message: "Areas retrieved successfully",
		Data:    areas,
	})
}
