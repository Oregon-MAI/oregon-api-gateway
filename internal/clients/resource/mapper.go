package resource

import (
	"time"

	resourcev1 "github.com/acyushka/oregon-infra/contracts/gen/go/resource"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type GetAvailableResourcesDTO struct {
	Resources  []ResourceDTO `json:"resources"`
	TotalCount int32         `json:"total_count"`
}

type ResourceDTO struct {
	ResourceID string    `json:"resource_id"`
	Name       string    `json:"name"`
	Type       string    `json:"type"`
	Location   string    `json:"location"`
	Status     string    `json:"status"`
	Details    any       `json:"details,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
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

func FromGetAvailableResourcesResponse(resp *resourcev1.GetAvailableResourcesResponse) *GetAvailableResourcesDTO {
	if resp == nil {
		return nil
	}

	dto := &GetAvailableResourcesDTO{
		TotalCount: resp.TotalCount,
		Resources:  make([]ResourceDTO, 0, len(resp.Resources)),
	}

	for _, r := range resp.Resources {
		dto.Resources = append(dto.Resources, *FromResource(r))
	}

	return dto
}

func FromResource(r *resourcev1.Resource) *ResourceDTO {
	if r == nil {
		return nil
	}

	dto := &ResourceDTO{
		ResourceID: r.ResourceId,
		Name:       r.Name,
		Type:       r.Type.String(),
		Location:   r.Location,
		Status:     r.Status.String(),
		CreatedAt:  fromTimestamp(r.CreatedAt),
		UpdatedAt:  fromTimestamp(r.UpdatedAt),
	}

	switch details := r.Details.(type) {
	case *resourcev1.Resource_MeetingRoom:
		dto.Details = MeetingRoomDTO{
			Capacity:      details.MeetingRoom.Capacity,
			HasProjector:  details.MeetingRoom.HasProjector,
			HasWhiteboard: details.MeetingRoom.HasWhiteboard,
		}
	case *resourcev1.Resource_Workspace:
		dto.Details = WorkspaceDTO{
			HasMonitor: details.Workspace.HasMonitor,
		}
	case *resourcev1.Resource_Device:
		dto.Details = DeviceDTO{
			DeviceType:   details.Device.DeviceType,
			SerialNumber: details.Device.SerialNumber,
			Model:        details.Device.Model,
			Description:  details.Device.Description,
		}
	}

	return dto
}

func FromCheckResourceStatusResponse(resp *resourcev1.CheckResourceStatusResponse) *CheckResourceStatusDTO {
	if resp == nil {
		return nil
	}
	return &CheckResourceStatusDTO{
		IsAvailable: resp.IsAvailable,
		Status:      resp.Status.String(),
	}
}

func ToGetAvailableResourcesRequest(types []string, location string) *resourcev1.GetAvailableResourcesRequest {
	req := &resourcev1.GetAvailableResourcesRequest{
		Location: location,
		Types:    make([]resourcev1.ResourceType, 0, len(types)),
	}

	for _, t := range types {
		if rt, ok := resourcev1.ResourceType_value[t]; ok {
			req.Types = append(req.Types, resourcev1.ResourceType(rt))
		}
	}

	return req
}

func ToGetResourceRequest(resourceID string) *resourcev1.GetResourceRequest {
	return &resourcev1.GetResourceRequest{
		ResourceId: resourceID,
	}
}

func ToCheckResourceStatusRequest(resourceID string) *resourcev1.CheckResourceStatusRequest {
	return &resourcev1.CheckResourceStatusRequest{
		ResourceId: resourceID,
	}
}

func fromTimestamp(ts *timestamppb.Timestamp) time.Time {
	if ts == nil {
		return time.Time{}
	}
	return ts.AsTime()
}

func toTimestamp(t time.Time) *timestamppb.Timestamp {
	return timestamppb.New(t)
}
