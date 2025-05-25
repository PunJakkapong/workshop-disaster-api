package models

// Assignment model
type Assignment struct {
	AreaID             string         `json:"area_id"`
	TruckID            string         `json:"truck_id"`
	ResourcesDelivered map[string]int `json:"resources_delivered"`
	Message            string         `json:"message,omitempty"`
}

// // AssignmentResponse represents the response for truck assignments
// type AssignmentResponse struct {
// 	Assignments []AreaAssignment `json:"assignments"`
// }

// // AssignmentResult represents the assignment for a specific area
// type AssignmentResult struct {
// 	AreaID             string         `json:"area_id"`
// 	TruckID            string         `json:"truck_id,omitempty"`
// 	ResourcesDelivered map[string]int `json:"resources_delivered,omitempty"`
// 	Message            string         `json:"message,omitempty"`
// }
