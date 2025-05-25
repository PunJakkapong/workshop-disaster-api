package models

// Assignment model
type Assignment struct {
	AreaID             string         `json:"area_id"`
	TruckID            string         `json:"truck_id"`
	ResourcesDelivered map[string]int `json:"resources_delivered"`
}

// AssignmentRequest represents the request body for truck assignments
type AssignmentRequest struct {
	Areas []AreaAssignmentRequest `json:"areas" binding:"required"`
}

// AreaAssignmentRequest represents an area's assignment request
type AreaAssignmentRequest struct {
	AreaID           string         `json:"areaId" binding:"required"`
	RequiredResource map[string]int `json:"requiredResource" binding:"required"`
	Urgency          int            `json:"urgency" binding:"required"`
	TimeConstraint   int            `json:"timeConstraint" binding:"required"` // in minutes
}

// AssignmentResponse represents the response for truck assignments
type AssignmentResponse struct {
	Assignments []AreaAssignment `json:"assignments"`
}

// AreaAssignment represents the assignment for a specific area
type AreaAssignment struct {
	AreaID     string   `json:"areaId"`
	TruckIDs   []string `json:"truckIds"`
	TotalTime  int      `json:"totalTime"` // in minutes
	TotalCost  float64  `json:"totalCost"`
	IsComplete bool     `json:"isComplete"` // whether all required resources are met
}
