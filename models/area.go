package models

// Area for get areas
type Area struct {
	AreaID            string         `json:"areaId"`
	UrgencyLevel      int            `json:"urgencyLevel"`
	RequiredResources map[string]int `json:"requiredResources"`
	TimeConstraint    int            `json:"timeConstraint"`
}

// CreateAreaRequest for create area
type CreateAreaRequest struct {
	AreaID            string         `json:"areaId" binding:"required"`
	UrgencyLevel      int            `json:"urgencyLevel" binding:"required,min=1,max=5"`
	RequiredResources map[string]int `json:"requiredResources" binding:"required,dive,min=0"`
	TimeConstraint    int            `json:"timeConstraint" binding:"required,min=0"`
}
