package controllers

import (
	"database/sql"
	"encoding/json"
	"math"
	"net/http"
	"sort"
	"workship-disaster-api/models"
	"workship-disaster-api/resp"

	"github.com/gin-gonic/gin"
)

type AssignmentController struct {
	db *sql.DB
}

func NewAssignmentController(db *sql.DB) *AssignmentController {
	return &AssignmentController{db: db}
}

// CreateAssignment handles the creation of truck assignments for areas
func (c *AssignmentController) CreateAssignment(ctx *gin.Context) {
	var req models.AssignmentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, resp.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Binding error",
			Error:   err.Error(),
		})
		return
	}

	// Get all trucks
	rows, err := c.db.Query("SELECT truck_id, available_resources, travel_time_to_area FROM trucks")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, resp.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to fetch trucks",
			Error:   err.Error(),
		})
		return
	}
	defer rows.Close()

	// Parse truck data
	type TruckData struct {
		ID                 string
		AvailableResources map[string]int
		TravelTimeToArea   map[string]int
	}

	var trucks []TruckData
	for rows.Next() {
		var truck TruckData
		var resourcesJSON, travelTimeJSON []byte
		if err := rows.Scan(&truck.ID, &resourcesJSON, &travelTimeJSON); err != nil {
			ctx.JSON(http.StatusInternalServerError, resp.ErrorResponse{
				Code:    http.StatusInternalServerError,
				Message: "Failed to parse truck data",
				Error:   err.Error(),
			})
			return
		}

		if err := json.Unmarshal(resourcesJSON, &truck.AvailableResources); err != nil {
			ctx.JSON(http.StatusInternalServerError, resp.ErrorResponse{
				Code:    http.StatusInternalServerError,
				Message: "Failed to parse truck resources",
				Error:   err.Error(),
			})
			return
		}

		if err := json.Unmarshal(travelTimeJSON, &truck.TravelTimeToArea); err != nil {
			ctx.JSON(http.StatusInternalServerError, resp.ErrorResponse{
				Code:    http.StatusInternalServerError,
				Message: "Failed to parse truck travel times",
				Error:   err.Error(),
			})
			return
		}

		trucks = append(trucks, truck)
	}

	// Sort areas by urgency (descending)
	sort.Slice(req.Areas, func(i, j int) bool {
		return req.Areas[i].Urgency > req.Areas[j].Urgency
	})

	// Process assignments
	assignments := make([]models.AreaAssignment, 0)
	assignedTrucks := make(map[string]bool)

	for _, area := range req.Areas {
		assignment := models.AreaAssignment{
			AreaID:     area.AreaID,
			TruckIDs:   make([]string, 0),
			TotalTime:  0,
			TotalCost:  0,
			IsComplete: false,
		}

		remainingResources := make(map[string]int)
		for resource, amount := range area.RequiredResource {
			remainingResources[resource] = amount
		}

		// Find suitable trucks for this area
		for _, truck := range trucks {
			if assignedTrucks[truck.ID] {
				continue
			}

			travelTime, exists := truck.TravelTimeToArea[area.AreaID]
			if !exists || travelTime > area.TimeConstraint {
				continue
			}

			// Check if truck has any needed resources
			hasNeededResources := false
			for resource := range remainingResources {
				if truck.AvailableResources[resource] > 0 {
					hasNeededResources = true
					break
				}
			}

			if !hasNeededResources {
				continue
			}

			// Assign truck to area
			assignment.TruckIDs = append(assignment.TruckIDs, truck.ID)
			assignedTrucks[truck.ID] = true
			assignment.TotalTime = travelTime

			// Update remaining resources
			for resource, amount := range remainingResources {
				available := truck.AvailableResources[resource]
				if available > 0 {
					remainingResources[resource] = int(math.Max(0, float64(amount-available)))
				}
			}

			// Calculate cost (example: cost = time * urgency)
			assignment.TotalCost += float64(travelTime) * float64(area.Urgency)

			// Check if all resources are met
			allResourcesMet := true
			for _, amount := range remainingResources {
				if amount > 0 {
					allResourcesMet = false
					break
				}
			}
			assignment.IsComplete = allResourcesMet

			if allResourcesMet {
				break
			}
		}

		assignments = append(assignments, assignment)
	}

	ctx.JSON(http.StatusOK, resp.SuccessResponse{
		Code:    http.StatusOK,
		Message: "Assignments created successfully",
		Data: models.AssignmentResponse{
			Assignments: assignments,
		},
	})
}
