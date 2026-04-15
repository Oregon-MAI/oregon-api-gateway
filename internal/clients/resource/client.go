package resource

import (
	"log/slog"

	"github.com/OnYyon/oregon-api-gateway/internal/clients/grpc"
	resourcev1 "github.com/acyushka/oregon-infra/contracts/gen/go/resource"
)

type Client struct {
	publicGrpcClient  *grpc.Client
	bookingGrpcClient *grpc.Client
	publicClient      resourcev1.ResourcePublicServiceClient
	bookingClient     resourcev1.ResourceBookingServiceClient
	adminClient       resourcev1.ResourcePublicServiceClient
	log               *slog.Logger
}

func NewClient(publicCfg, bookingCfg *grpc.Config, log *slog.Logger) (*Client, error) {
	publicGrpcClient, err := grpc.NewGRPCClient(publicCfg, log)
	if err != nil {
		return nil, err
	}

	bookingGrpcClient, err := grpc.NewGRPCClient(bookingCfg, log)
	if err != nil {
		if closeErr := publicGrpcClient.Close(); closeErr != nil {
			log.Error("failed to close public grpc client", slog.Any("error", closeErr))
		}
		return nil, err
	}

	return &Client{
		publicGrpcClient:  publicGrpcClient,
		bookingGrpcClient: bookingGrpcClient,
		publicClient:      resourcev1.NewResourcePublicServiceClient(publicGrpcClient.Conn()),
		bookingClient:     resourcev1.NewResourceBookingServiceClient(bookingGrpcClient.Conn()),
		adminClient:       resourcev1.NewResourcePublicServiceClient(publicGrpcClient.Conn()),
		log:               log.With(slog.String("component", "resource_client")),
	}, nil
}

func (c *Client) PublicGRPCClient() *grpc.Client {
	return c.publicGrpcClient
}

func (c *Client) BookingGRPCClient() *grpc.Client {
	return c.bookingGrpcClient
}

func (c *Client) Close() error {
	err1 := c.publicGrpcClient.Close()
	err2 := c.bookingGrpcClient.Close()
	if err1 != nil {
		return err1
	}
	return err2
}

func (c *Client) PublicClient() resourcev1.ResourcePublicServiceClient {
	return c.publicClient
}

func (c *Client) BookingClient() resourcev1.ResourceBookingServiceClient {
	return c.bookingClient
}

func (c *Client) AdminClient() resourcev1.ResourcePublicServiceClient {
	return c.adminClient
}
