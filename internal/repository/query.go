package repository

// ResourceQuery represents resource query parameters
type ResourceQuery struct {
	Name     string `json:"name,omitempty"`
	Type     string `json:"type,omitempty"`
	ParentID *uint  `json:"parent_id,omitempty"`
	Status   int    `json:"status,omitempty"`
	Page     int    `json:"page,omitempty"`
	PageSize int    `json:"page_size,omitempty"`
}

// ActionQuery represents action query parameters
type ActionQuery struct {
	Name     string `json:"name,omitempty"`
	Category string `json:"category,omitempty"`
	Page     int    `json:"page,omitempty"`
	PageSize int    `json:"page_size,omitempty"`
}