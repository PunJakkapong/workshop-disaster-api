package models

// Truck rfor get trucks
type Truck struct {
	TruckID            string         `json:"truckId"`
	AvailableResources map[string]int `json:"availableResources"`
	TravelTimeToArea   map[string]int `json:"travelTimeToArea"`
}

// CreateTruckRequest for create truck
type CreateTruckRequest struct {
	TruckID            string         `json:"truckId" binding:"required"`
	AvailableResources map[string]int `json:"availableResources" binding:"required,dive,min=0"`
	TravelTimeToArea   map[string]int `json:"travelTimeToArea" binding:"required,dive,min=0"`
}
