package service

import (
	"fmt"
	"strings"
	"workship-disaster-api/models"
)

type AssignmentService struct {
	areaService  *AreaService
	truckService *TruckService
}

func NewAssignmentService(areaService *AreaService, truckService *TruckService) *AssignmentService {
	return &AssignmentService{areaService, truckService}
}

// CreateAssignments use for api assignments
func (s *AssignmentService) CreateAssignments() ([]models.Assignment, error) {
	areas, err := s.areaService.GetAllAreas()
	if err != nil {
		return nil, fmt.Errorf("failed to get areas: %w", err)
	}

	trucks, err := s.truckService.GetAllTrucks()
	if err != nil {
		return nil, fmt.Errorf("failed to get trucks: %w", err)
	}

	// Keep track of used trucks
	usedTrucks := make(map[string]bool)
	assignments := []models.Assignment{}

	for _, area := range areas {
		var bestTruck *TruckData
		var bestTravelTime int
		var deliveredResources map[string]int

		for _, truck := range trucks {
			if usedTrucks[truck.ID] {
				continue
			}

			travelTime, ok := truck.TravelTimeToArea[area.ID]
			// If truck travel time not enough
			if !ok || travelTime > area.TimeConstraint {
				continue
			}

			// If truck resource can fill area required resource
			if canFulfill(truck.AvailableResources, area.RequiredResource) {
				// If current truck travel time is best
				if bestTruck == nil || travelTime < bestTravelTime {
					bestTruck = &truck
					bestTravelTime = travelTime
				}
			}
		}

		// If founded the best truck
		if bestTruck != nil {
			assignedTruckID := bestTruck.ID
			deliveredResources = area.RequiredResource
			usedTrucks[assignedTruckID] = true

			assignments = append(assignments, models.Assignment{
				AreaID:             area.ID,
				TruckID:            assignedTruckID,
				ResourcesDelivered: deliveredResources,
			})
		} else {
			msg := "No trucks available, Remaining resources needed:"
			parts := []string{}
			for k, v := range area.RequiredResource {
				parts = append(parts, fmt.Sprintf(" %s: %d", k, v))
			}
			assignments = append(assignments, models.Assignment{
				AreaID:  area.ID,
				Message: msg + strings.Join(parts, ","),
			})
		}
	}

	return assignments, nil
}

func canFulfill(available, required map[string]int) bool {
	for k, v := range required {
		if available[k] < v {
			return false
		}
	}
	return true
}
