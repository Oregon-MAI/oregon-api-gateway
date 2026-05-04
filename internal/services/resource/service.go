package resource

import (
	"context"
	"fmt"

	"github.com/OnYyon/oregon-api-gateway/internal/clients/grpc"
	"github.com/OnYyon/oregon-api-gateway/internal/clients/resource"
	resourcev1 "github.com/Oregon-MAI/oregon-infrastructure/contracts/gen/go/resource"
)

type Service struct {
	client *resource.Client
}

func NewService(client *resource.Client) *Service {
	return &Service{client: client}
}

func (s *Service) GetAvailableResources(ctx context.Context, types []string, location string) (*resourcev1.GetAvailableResourcesResponse, error) {
	req := &resourcev1.GetAvailableResourcesRequest{
		Location: location,
		Types:    make([]resourcev1.ResourceType, 0, len(types)),
	}

	for _, t := range types {
		if rt, ok := resourcev1.ResourceType_value[t]; ok {
			req.Types = append(req.Types, resourcev1.ResourceType(rt))
		}
	}

	resp, err := grpc.Call(
		ctx,
		s.client.PublicGRPCClient().Conn(),
		s.client.PublicGRPCClient().Log(),
		s.client.PublicGRPCClient().Tracer(),
		s.client.PublicGRPCClient().Timeout(),
		"Resource.GetAvailableResources",
		func(ctx context.Context, r *resourcev1.GetAvailableResourcesRequest) (*resourcev1.GetAvailableResourcesResponse, error) {
			return s.client.PublicClient().GetAvailableResources(ctx, r)
		}, req)

	if err != nil {
		return nil, resource.MapGRPCErr(err)
	}

	return resp, nil
}

func (s *Service) GetResource(ctx context.Context, resourceID string) (*resourcev1.GetResourceResponse, error) {
	req := &resourcev1.GetResourceRequest{ResourceId: resourceID}

	resp, err := grpc.Call(
		ctx,
		s.client.PublicGRPCClient().Conn(),
		s.client.PublicGRPCClient().Log(),
		s.client.PublicGRPCClient().Tracer(),
		s.client.PublicGRPCClient().Timeout(),
		"Resource.GetResource",
		func(ctx context.Context, r *resourcev1.GetResourceRequest) (*resourcev1.GetResourceResponse, error) {
			return s.client.PublicClient().GetResource(ctx, r)
		}, req)

	if err != nil {
		return nil, resource.MapGRPCErr(err)
	}

	return resp, nil
}

func (s *Service) CheckResourceStatus(ctx context.Context, resourceID string) (*resourcev1.CheckResourceStatusResponse, error) {
	req := &resourcev1.CheckResourceStatusRequest{ResourceId: resourceID}

	resp, err := grpc.Call(
		ctx,
		s.client.BookingGRPCClient().Conn(),
		s.client.BookingGRPCClient().Log(),
		s.client.BookingGRPCClient().Tracer(),
		s.client.BookingGRPCClient().Timeout(),
		"Resource.CheckResourceStatus",
		func(ctx context.Context, r *resourcev1.CheckResourceStatusRequest) (*resourcev1.CheckResourceStatusResponse, error) {
			return s.client.BookingClient().CheckResourceStatus(ctx, r)
		}, req)

	if err != nil {
		return nil, resource.MapGRPCErr(err)
	}

	return resp, nil
}

func (s *Service) CreateResource(ctx context.Context, req *resourcev1.CreateResourceRequest) (*resourcev1.CreateResourceResponse, error) {
	resp, err := grpc.Call(
		ctx,
		s.client.PublicGRPCClient().Conn(),
		s.client.PublicGRPCClient().Log(),
		s.client.PublicGRPCClient().Tracer(),
		s.client.PublicGRPCClient().Timeout(),
		"Resource.CreateResource",
		func(ctx context.Context, r *resourcev1.CreateResourceRequest) (*resourcev1.CreateResourceResponse, error) {
			return s.client.AdminClient().CreateResource(ctx, r)
		}, req)

	if err != nil {
		return nil, resource.MapGRPCErr(err)
	}

	return resp, nil
}

func (s *Service) GetResourcesList(ctx context.Context, types []string) (*resourcev1.GetResourcesListResponse, error) {
	req := &resourcev1.GetResourcesListRequest{
		Types: make([]resourcev1.ResourceType, 0, len(types)),
	}

	for _, t := range types {
		if rt, ok := resourcev1.ResourceType_value[t]; ok {
			req.Types = append(req.Types, resourcev1.ResourceType(rt))
		}
	}

	resp, err := grpc.Call(
		ctx,
		s.client.PublicGRPCClient().Conn(),
		s.client.PublicGRPCClient().Log(),
		s.client.PublicGRPCClient().Tracer(),
		s.client.PublicGRPCClient().Timeout(),
		"Resource.GetResourcesList",
		func(ctx context.Context, r *resourcev1.GetResourcesListRequest) (*resourcev1.GetResourcesListResponse, error) {
			return s.client.AdminClient().GetResourcesList(ctx, r)
		}, req)

	if err != nil {
		return nil, resource.MapGRPCErr(err)
	}

	return resp, nil
}

func (s *Service) UpdateResource(ctx context.Context, dto *resource.UpdateResourceRequestDTO) (*resource.ResourceDTO, error) {
	req := resource.ToUpdateResourceRequest(dto)

	resp, err := grpc.Call(
		ctx,
		s.client.PublicGRPCClient().Conn(),
		s.client.PublicGRPCClient().Log(),
		s.client.PublicGRPCClient().Tracer(),
		s.client.PublicGRPCClient().Timeout(),
		"Resource.UpdateResource",
		func(ctx context.Context, r *resourcev1.UpdateResourceRequest) (*resourcev1.UpdateResourceResponse, error) {
			return s.client.AdminClient().UpdateResource(ctx, r)
		}, req)
	fmt.Println(req, "\n", resp)
	if err != nil {
		return nil, resource.MapGRPCErr(err)
	}

	return resource.FromResource(resp.Resource), nil
}

func (s *Service) DeleteResource(ctx context.Context, resourceID string) (bool, error) {
	req := &resourcev1.DeleteResourceRequest{ResourceId: resourceID}

	resp, err := grpc.Call(
		ctx,
		s.client.PublicGRPCClient().Conn(),
		s.client.PublicGRPCClient().Log(),
		s.client.PublicGRPCClient().Tracer(),
		s.client.PublicGRPCClient().Timeout(),
		"Resource.DeleteResource",
		func(ctx context.Context, r *resourcev1.DeleteResourceRequest) (*resourcev1.DeleteResourceResponse, error) {
			return s.client.AdminClient().DeleteResource(ctx, r)
		}, req)

	if err != nil {
		return false, resource.MapGRPCErr(err)
	}

	return resp.Success, nil
}

func (s *Service) ChangeResourceStatus(ctx context.Context, req *resourcev1.ChangeResourceStatusRequest) (*resourcev1.ChangeResourceStatusResponse, error) {
	resp, err := grpc.Call(
		ctx,
		s.client.PublicGRPCClient().Conn(),
		s.client.PublicGRPCClient().Log(),
		s.client.PublicGRPCClient().Tracer(),
		s.client.PublicGRPCClient().Timeout(),
		"Resource.ChangeResourceStatus",
		func(ctx context.Context, r *resourcev1.ChangeResourceStatusRequest) (*resourcev1.ChangeResourceStatusResponse, error) {
			return s.client.AdminClient().ChangeResourceStatus(ctx, r)
		}, req)

	if err != nil {
		return nil, resource.MapGRPCErr(err)
	}

	return resp, nil
}

func (s *Service) UpdateResourceOccupancy(ctx context.Context, req *resourcev1.UpdateResourceOccupancyRequest) (*resourcev1.UpdateResourceOccupancyResponse, error) {
	resp, err := grpc.Call(
		ctx,
		s.client.BookingGRPCClient().Conn(),
		s.client.BookingGRPCClient().Log(),
		s.client.BookingGRPCClient().Tracer(),
		s.client.BookingGRPCClient().Timeout(),
		"ResourceBooking.UpdateResourceOccupancy",
		func(ctx context.Context, r *resourcev1.UpdateResourceOccupancyRequest) (*resourcev1.UpdateResourceOccupancyResponse, error) {
			return s.client.BookingClient().UpdateResourceOccupancy(ctx, r)
		}, req)

	if err != nil {
		return nil, resource.MapGRPCErr(err)
	}

	return resp, nil
}
