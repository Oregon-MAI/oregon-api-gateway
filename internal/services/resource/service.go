package resource

import (
	"context"

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

func (s *Service) GetAvailableResources(ctx context.Context, types []string, location string) (*resource.GetAvailableResourcesDTO, error) {
	req := resource.ToGetAvailableResourcesRequest(types, location)

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

	return resource.FromGetAvailableResourcesResponse(resp), nil
}

func (s *Service) GetResource(ctx context.Context, resourceID string) (*resource.ResourceDTO, error) {
	req := resource.ToGetResourceRequest(resourceID)

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

	return resource.FromResource(resp.Resource), nil
}

func (s *Service) CheckResourceStatus(ctx context.Context, resourceID string) (*resource.CheckResourceStatusDTO, error) {
	req := resource.ToCheckResourceStatusRequest(resourceID)

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

	return resource.FromCheckResourceStatusResponse(resp), nil
}

func (s *Service) CreateResource(ctx context.Context, dto *resource.CreateResourceRequestDTO) (*resource.ResourceDTO, error) {
	req := resource.ToCreateResourceRequest(dto)

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

	return resource.FromResource(resp.Resource), nil
}

func (s *Service) GetResourcesList(ctx context.Context, types []string) (*resource.GetResourcesListDTO, error) {
	req := resource.ToGetResourcesListRequest(types)

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

	return resource.FromGetResourcesListResponse(resp), nil
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

	if err != nil {
		return nil, resource.MapGRPCErr(err)
	}

	return resource.FromResource(resp.Resource), nil
}

func (s *Service) DeleteResource(ctx context.Context, resourceID string) (bool, error) {
	req := resource.ToDeleteResourceRequest(resourceID)

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

func (s *Service) ChangeResourceStatus(ctx context.Context, dto *resource.ChangeResourceStatusRequestDTO) (*resource.ResourceDTO, error) {
	req := resource.ToChangeResourceStatusRequest(dto)

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

	return resource.FromResource(resp.Resource), nil
}

func (s *Service) UpdateResourceOccupancy(ctx context.Context, dto *resource.UpdateResourceOccupancyRequestDTO) (*resource.UpdateResourceOccupancyResponseDTO, error) {
	req := resource.ToUpdateResourceOccupancyRequest(dto)

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

	return resource.FromUpdateResourceOccupancyResponse(resp), nil
}
