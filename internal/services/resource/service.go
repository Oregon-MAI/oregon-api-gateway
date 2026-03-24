package resource

import (
	"context"

	"github.com/OnYyon/oregon-api-gateway/internal/clients/grpc"
	"github.com/OnYyon/oregon-api-gateway/internal/clients/resource"
	resourcev1 "github.com/acyushka/oregon-infra/contracts/gen/go/resource"
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
		s.client.GRPCClient().Conn(),
		s.client.GRPCClient().Log(),
		s.client.GRPCClient().Tracer(),
		s.client.GRPCClient().Timeout(),
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
		s.client.GRPCClient().Conn(),
		s.client.GRPCClient().Log(),
		s.client.GRPCClient().Tracer(),
		s.client.GRPCClient().Timeout(),
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
		s.client.GRPCClient().Conn(),
		s.client.GRPCClient().Log(),
		s.client.GRPCClient().Tracer(),
		s.client.GRPCClient().Timeout(),
		"Resource.CheckResourceStatus",
		func(ctx context.Context, r *resourcev1.CheckResourceStatusRequest) (*resourcev1.CheckResourceStatusResponse, error) {
			return s.client.BookingClient().CheckResourceStatus(ctx, r)
		}, req)

	if err != nil {
		return nil, resource.MapGRPCErr(err)
	}

	return resource.FromCheckResourceStatusResponse(resp), nil
}
