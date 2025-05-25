package service

import (
	"fmt"
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

		// Variables to handle edge cases
		hasTruckWithTravelTimeEntry := false
		hasTruckWithSufficientResources := false
		hasTruckWithSufficientResourcesAndTime := false

		for _, truck := range trucks {
			if usedTrucks[truck.ID] {
				continue
			}

			travelTime, ok := truck.TravelTimeToArea[area.ID]
			if ok {
				hasTruckWithTravelTimeEntry = true
			}

			// If truck resource can fill area required resource
			if canFulfill(truck.AvailableResources, area.RequiredResource) {
				hasTruckWithSufficientResources = true

				if ok && travelTime <= area.TimeConstraint {
					hasTruckWithSufficientResourcesAndTime = true

					// If current truck travel time is best
					if bestTruck == nil || travelTime < bestTravelTime {
						bestTruck = &truck
						bestTravelTime = travelTime
					}
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
			// Create detailed fallback message
			var msg string
			if !hasTruckWithTravelTimeEntry {
				msg = "No trucks have a valid route to this area."
			} else if !hasTruckWithSufficientResources {
				msg = "No truck has sufficient resources to fulfill this area's needs."
			} else if !hasTruckWithSufficientResourcesAndTime {
				msg = "All trucks with sufficient resources exceed the time constraint."
			} else {
				msg = "No trucks available for assignment."
			}

			assignments = append(assignments, models.Assignment{
				AreaID:  area.ID,
				Message: msg,
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
