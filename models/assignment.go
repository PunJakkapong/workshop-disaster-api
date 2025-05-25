package models

// Assignment model
type Assignment struct {
	AreaID             string         `json:"area_id"`
	TruckID            string         `json:"truck_id"`
	ResourcesDelivered map[string]int `json:"resources_delivered"`
	Message            string         `json:"message,omitempty"`
}
