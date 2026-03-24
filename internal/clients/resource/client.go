package resource

import (
	"log/slog"

	"github.com/OnYyon/oregon-api-gateway/internal/clients/grpc"
	resourcev1 "github.com/acyushka/oregon-infra/contracts/gen/go/resource"
)

type Client struct {
	grpcClient    *grpc.Client
	publicClient  resourcev1.ResourcePublicServiceClient
	bookingClient resourcev1.ResourceBookingServiceClient
	adminClient   resourcev1.ResourceAdminServiceClient
	log           *slog.Logger
}

func NewClient(cfg *grpc.Config, log *slog.Logger) (*Client, error) {
	grpcClient, err := grpc.NewGRPCClient(*cfg, log)
	if err != nil {
		return nil, err
	}

	return &Client{
		grpcClient:    grpcClient,
		publicClient:  resourcev1.NewResourcePublicServiceClient(grpcClient.Conn()),
		bookingClient: resourcev1.NewResourceBookingServiceClient(grpcClient.Conn()),
		adminClient:   resourcev1.NewResourceAdminServiceClient(grpcClient.Conn()),
		log:           log.With(slog.String("component", "resource_client")),
	}, nil
}

func (c *Client) GRPCClient() *grpc.Client {
	return c.grpcClient
}

func (c *Client) Close() error {
	return c.grpcClient.Close()
}

func (c *Client) PublicClient() resourcev1.ResourcePublicServiceClient {
	return c.publicClient
}

func (c *Client) BookingClient() resourcev1.ResourceBookingServiceClient {
	return c.bookingClient
}

func (c *Client) AdminClient() resourcev1.ResourceAdminServiceClient {
	return c.adminClient
}
