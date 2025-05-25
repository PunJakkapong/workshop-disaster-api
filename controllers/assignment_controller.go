package controllers

import (
	"database/sql"
	"net/http"
	"workship-disaster-api/resp"
	"workship-disaster-api/service"

	"github.com/gin-gonic/gin"
)

// AssignmentController ...
type AssignmentController struct {
	db                *sql.DB
	assignmentService *service.AssignmentService
}

// NewAssignmentController ...
func NewAssignmentController(db *sql.DB) *AssignmentController {
	areaService := service.NewAreaService(db)
	truckService := service.NewTruckService(db)
	assignmentService := service.NewAssignmentService(areaService, truckService)

	return &AssignmentController{
		db:                db,
		assignmentService: assignmentService,
	}
}

// CreateAssignment handles the creation of truck assignments for areas
func (c *AssignmentController) CreateAssignment(ctx *gin.Context) {
	assignments, err := c.assignmentService.CreateAssignments()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, resp.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to create assignments",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, resp.SuccessResponse{
		Code:    http.StatusOK,
		Message: "Assignments created successfully",
		Data:    assignments,
	})
}
