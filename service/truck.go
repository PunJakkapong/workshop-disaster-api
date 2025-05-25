package service

import (
	"database/sql"
	"encoding/json"
	"fmt"
)

type TruckData struct {
	ID                 string
	AvailableResources map[string]int
	TravelTimeToArea   map[string]int
}

type TruckService struct {
	db *sql.DB
}

func NewTruckService(db *sql.DB) *TruckService {
	return &TruckService{db: db}
}

// GetAllTrucks fetches all trucks from the database
func (s *TruckService) GetAllTrucks() ([]TruckData, error) {
	truckRows, err := s.db.Query("SELECT truck_id, available_resources, travel_time_to_area FROM trucks")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch trucks: %w", err)
	}
	defer truckRows.Close()

	var trucks []TruckData
	for truckRows.Next() {
		var truck TruckData
		var resourcesJSON, travelTimeJSON []byte
		if err := truckRows.Scan(&truck.ID, &resourcesJSON, &travelTimeJSON); err != nil {
			return nil, fmt.Errorf("failed to parse truck data: %w", err)
		}

		if err := json.Unmarshal(resourcesJSON, &truck.AvailableResources); err != nil {
			return nil, fmt.Errorf("failed to parse truck resources: %w", err)
		}

		if err := json.Unmarshal(travelTimeJSON, &truck.TravelTimeToArea); err != nil {
			return nil, fmt.Errorf("failed to parse truck travel times: %w", err)
		}

		trucks = append(trucks, truck)
	}

	return trucks, nil
}
