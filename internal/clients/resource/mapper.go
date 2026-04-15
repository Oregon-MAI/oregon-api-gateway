package resource

import (
	"encoding/json"
	"time"

	resourcev1 "github.com/acyushka/oregon-infra/contracts/gen/go/resource"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

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
		if b, err := json.Marshal(MeetingRoomDTO{
			Capacity:      details.MeetingRoom.Capacity,
			HasProjector:  details.MeetingRoom.HasProjector,
			HasWhiteboard: details.MeetingRoom.HasWhiteboard,
		}); err == nil {
			dto.Details = b
		}
	case *resourcev1.Resource_Workspace:
		if b, err := json.Marshal(WorkspaceDTO{
			HasMonitor: details.Workspace.HasMonitor,
		}); err == nil {
			dto.Details = b
		}
	case *resourcev1.Resource_Device:
		if b, err := json.Marshal(DeviceDTO{
			DeviceType:   details.Device.DeviceType,
			SerialNumber: details.Device.SerialNumber,
			Model:        details.Device.Model,
			Description:  details.Device.Description,
		}); err == nil {
			dto.Details = b
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

func FromGetResourcesListResponse(resp *resourcev1.GetResourcesListResponse) *GetResourcesListDTO {
	if resp == nil {
		return nil
	}

	dto := &GetResourcesListDTO{
		Resources: make([]ResourceDTO, 0, len(resp.Resources)),
	}

	for _, r := range resp.Resources {
		dto.Resources = append(dto.Resources, *FromResource(r))
	}

	return dto
}

func FromUpdateResourceOccupancyResponse(resp *resourcev1.UpdateResourceOccupancyResponse) *UpdateResourceOccupancyResponseDTO {
	if resp == nil {
		return nil
	}
	return &UpdateResourceOccupancyResponseDTO{
		Success: resp.Success,
		Status:  resp.Status.String(),
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

func ToCreateResourceRequest(req *CreateResourceRequestDTO) *resourcev1.CreateResourceRequest {
	if req == nil {
		return nil
	}

	pbReq := &resourcev1.CreateResourceRequest{
		Name:     req.Name,
		Location: req.Location,
	}

	if rt, ok := resourcev1.ResourceType_value[req.Type]; ok {
		pbReq.Type = resourcev1.ResourceType(rt)
	}

	switch req.Type {
	case "RESOURCE_TYPE_MEETING_ROOM":
		var d MeetingRoomDTO
		if err := json.Unmarshal(req.Details, &d); err == nil {
			pbReq.Details = &resourcev1.CreateResourceRequest_MeetingRoom{
				MeetingRoom: &resourcev1.MeetingRoomDetails{
					Capacity:      d.Capacity,
					HasProjector:  d.HasProjector,
					HasWhiteboard: d.HasWhiteboard,
				},
			}
		}
	case "RESOURCE_TYPE_WORKSPACE":
		var d WorkspaceDTO
		if err := json.Unmarshal(req.Details, &d); err == nil {
			pbReq.Details = &resourcev1.CreateResourceRequest_Workspace{
				Workspace: &resourcev1.WorkspaceDetails{
					HasMonitor: d.HasMonitor,
				},
			}
		}
	case "RESOURCE_TYPE_DEVICE":
		var d DeviceDTO
		if err := json.Unmarshal(req.Details, &d); err == nil {
			pbReq.Details = &resourcev1.CreateResourceRequest_Device{
				Device: &resourcev1.DeviceDetails{
					DeviceType:   d.DeviceType,
					SerialNumber: d.SerialNumber,
					Model:        d.Model,
					Description:  d.Description,
				},
			}
		}
	}

	return pbReq
}

func ToGetResourcesListRequest(types []string) *resourcev1.GetResourcesListRequest {
	req := &resourcev1.GetResourcesListRequest{
		Types: make([]resourcev1.ResourceType, 0, len(types)),
	}

	for _, t := range types {
		if rt, ok := resourcev1.ResourceType_value[t]; ok {
			req.Types = append(req.Types, resourcev1.ResourceType(rt))
		}
	}

	return req
}

func ToUpdateResourceRequest(req *UpdateResourceRequestDTO) *resourcev1.UpdateResourceRequest {
	if req == nil {
		return nil
	}

	pbReq := &resourcev1.UpdateResourceRequest{
		ResourceId: req.ResourceID,
		Resource: &resourcev1.Resource{
			ResourceId: req.ResourceID,
			Name:       req.Name,
			Location:   req.Location,
		},
	}

	if rt, ok := resourcev1.ResourceType_value[req.Type]; ok {
		pbReq.Resource.Type = resourcev1.ResourceType(rt)
	}
	if rs, ok := resourcev1.ResourceStatus_value[req.Status]; ok {
		pbReq.Resource.Status = resourcev1.ResourceStatus(rs)
	}

	var detailsPath string

	switch req.Type {
	case "RESOURCE_TYPE_MEETING_ROOM":
		var d MeetingRoomDTO
		if err := json.Unmarshal(req.Details, &d); err == nil {
			pbReq.Resource.Details = &resourcev1.Resource_MeetingRoom{
				MeetingRoom: &resourcev1.MeetingRoomDetails{
					Capacity:      d.Capacity,
					HasProjector:  d.HasProjector,
					HasWhiteboard: d.HasWhiteboard,
				},
			}
			detailsPath = "meeting_room"
		}
	case "RESOURCE_TYPE_WORKSPACE":
		var d WorkspaceDTO
		if err := json.Unmarshal(req.Details, &d); err == nil {
			pbReq.Resource.Details = &resourcev1.Resource_Workspace{
				Workspace: &resourcev1.WorkspaceDetails{
					HasMonitor: d.HasMonitor,
				},
			}
			detailsPath = "workspace"
		}
	case "RESOURCE_TYPE_DEVICE":
		var d DeviceDTO
		if err := json.Unmarshal(req.Details, &d); err == nil {
			pbReq.Resource.Details = &resourcev1.Resource_Device{
				Device: &resourcev1.DeviceDetails{
					DeviceType:   d.DeviceType,
					SerialNumber: d.SerialNumber,
					Model:        d.Model,
					Description:  d.Description,
				},
			}
		}
		detailsPath = "device"
	}

	paths := []string{
		"name",
		"location",
		"type",
		"status",
		detailsPath,
	}
	pbReq.FieldMask = &fieldmaskpb.FieldMask{
		Paths: paths,
	}

	return pbReq
}

func ToDeleteResourceRequest(resourceID string) *resourcev1.DeleteResourceRequest {
	return &resourcev1.DeleteResourceRequest{
		ResourceId: resourceID,
	}
}

func ToUpdateResourceOccupancyRequest(req *UpdateResourceOccupancyRequestDTO) *resourcev1.UpdateResourceOccupancyRequest {
	if req == nil {
		return nil
	}
	return &resourcev1.UpdateResourceOccupancyRequest{
		ResourceId: req.ResourceID,
		IsOccupied: req.IsOccupied,
	}
}

func ToChangeResourceStatusRequest(req *ChangeResourceStatusRequestDTO) *resourcev1.ChangeResourceStatusRequest {
	if req == nil {
		return nil
	}

	pbReq := &resourcev1.ChangeResourceStatusRequest{
		ResourceId: req.ResourceID,
		Reason:     req.Reason,
	}

	if rs, ok := resourcev1.ResourceStatus_value[req.Status]; ok {
		pbReq.Status = resourcev1.ResourceStatus(rs)
	}

	return pbReq
}

func fromTimestamp(ts *timestamppb.Timestamp) time.Time {
	if ts == nil {
		return time.Time{}
	}
	return ts.AsTime()
}
