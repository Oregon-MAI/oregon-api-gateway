package resource

import (
	"encoding/json"
	"time"
)

type GetAvailableResourcesDTO struct {
	Resources  []ResourceDTO `json:"resources"`
	TotalCount int32         `json:"total_count"`
}

type ResourceDTO struct {
	ResourceID string          `json:"resource_id"`
	Name       string          `json:"name"`
	Type       string          `json:"type"`
	Location   string          `json:"location"`
	Status     string          `json:"status"`
	Details    json.RawMessage `json:"details,omitempty"`
	CreatedAt  time.Time       `json:"created_at"`
	UpdatedAt  time.Time       `json:"updated_at"`
}

type CreateResourceRequestDTO struct {
	Name     string          `json:"name"`
	Type     string          `json:"type"`
	Location string          `json:"location"`
	Details  json.RawMessage `json:"details,omitempty"`
}

type UpdateResourceRequestDTO struct {
	ResourceDTO
}

type ChangeResourceStatusRequestDTO struct {
	ResourceID string `json:"resource_id"`
	Status     string `json:"status"`
	Reason     string `json:"reason"`
}

type GetResourcesListDTO struct {
	Resources []ResourceDTO `json:"resources"`
}

type UpdateResourceOccupancyRequestDTO struct {
	ResourceID string `json:"resource_id"`
	IsOccupied bool   `json:"is_occupied"`
}

type UpdateResourceOccupancyResponseDTO struct {
	Success bool   `json:"success"`
	Status  string `json:"status"`
}

type MeetingRoomDTO struct {
	Capacity      int32 `json:"capacity"`
	HasProjector  bool  `json:"has_projector"`
	HasWhiteboard bool  `json:"has_whiteboard"`
}

type WorkspaceDTO struct {
	HasMonitor bool `json:"has_monitor"`
}

type DeviceDTO struct {
	DeviceType   string `json:"device_type"`
	SerialNumber string `json:"serial_number"`
	Model        string `json:"model"`
	Description  string `json:"description"`
}

type CheckResourceStatusDTO struct {
	IsAvailable bool   `json:"is_available"`
	Status      string `json:"status"`
}
