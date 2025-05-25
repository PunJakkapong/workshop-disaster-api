package service

import (
	"database/sql"
	"encoding/json"
	"fmt"
)

type AreaData struct {
	ID               string
	RequiredResource map[string]int
	Urgency          int
	TimeConstraint   int
}

type AreaService struct {
	db *sql.DB
}

func NewAreaService(db *sql.DB) *AreaService {
	return &AreaService{db: db}
}

// GetAllAreas fetches all areas from the database
func (s *AreaService) GetAllAreas() ([]AreaData, error) {
	areaRows, err := s.db.Query("SELECT area_id, required_resources, urgency_level, time_constraint FROM areas ORDER BY urgency_level DESC")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch areas: %w", err)
	}
	defer areaRows.Close()

	var areas []AreaData
	for areaRows.Next() {
		var area AreaData
		var resourcesJSON []byte
		if err := areaRows.Scan(&area.ID, &resourcesJSON, &area.Urgency, &area.TimeConstraint); err != nil {
			return nil, fmt.Errorf("failed to parse area data: %w", err)
		}

		if err := json.Unmarshal(resourcesJSON, &area.RequiredResource); err != nil {
			return nil, fmt.Errorf("failed to parse area resources: %w", err)
		}

		areas = append(areas, area)
	}

	return areas, nil
}
