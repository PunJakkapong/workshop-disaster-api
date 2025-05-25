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

type TruckController struct {
	db *sql.DB
}

func NewTruckController(db *sql.DB) *TruckController {
	return &TruckController{db: db}
}

// CreateTruck handles the creation of a new truck
func (c *TruckController) CreateTruck(ctx *gin.Context) {
	var req models.CreateTruckRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, resp.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Binding error.",
			Error:   err.Error(),
		})
		return
	}

	// Check if truck ID already exists
	var exists bool
	err := c.db.QueryRow("SELECT EXISTS(SELECT 1 FROM trucks WHERE truck_id = $1)", req.TruckID).Scan(&exists)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, resp.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to check truck ID",
			Error:   err.Error(),
		})
		return
	}
	if exists {
		ctx.JSON(http.StatusConflict, resp.ErrorResponse{
			Code:    http.StatusConflict,
			Message: "Truck ID already exists",
		})
		return
	}

	// Convert JSON fields
	resourcesJSON, err := json.Marshal(req.AvailableResources)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, resp.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to process available resources",
			Error:   err.Error(),
		})
		return
	}

	travelTimeJSON, err := json.Marshal(req.TravelTimeToArea)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, resp.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to process travel times",
			Error:   err.Error(),
		})
		return
	}

	// Insert into database
	_, err = c.db.Exec(
		"INSERT INTO trucks (truck_id, available_resources, travel_time_to_area) VALUES ($1, $2, $3)",
		req.TruckID, resourcesJSON, travelTimeJSON,
	)

	if err != nil {
		// Check if error is due to duplicate key
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			ctx.JSON(http.StatusConflict, resp.ErrorResponse{
				Code:    http.StatusConflict,
				Message: "Truck ID already exists",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, resp.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to create truck",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, resp.SuccessResponse{
		Code:    http.StatusCreated,
		Message: "Truck created successfully",
		Data: gin.H{
			"truckId": req.TruckID,
		},
	})
}
