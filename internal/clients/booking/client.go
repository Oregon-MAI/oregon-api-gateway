package booking

import (
	"log/slog"

	"github.com/OnYyon/oregon-api-gateway/internal/clients/grpc"
	bookingv1 "github.com/Oregon-MAI/oregon-infrastructure/contracts/gen/go/booking"
)

type Client struct {
	bookingGrpcClient *grpc.Client
	bookingClient     bookingv1.BookingServiceClient
	log               *slog.Logger
}

func NewClient(bookingCfg *grpc.Config, log *slog.Logger) (*Client, error) {
	bookingGrpcClient, err := grpc.NewGRPCClient(bookingCfg, log)
	if err != nil {
		return nil, err
	}

	return &Client{
		bookingGrpcClient: bookingGrpcClient,
		bookingClient:     bookingv1.NewBookingServiceClient(bookingGrpcClient.Conn()),
		log:               log.With(slog.String("component", "booking_client")),
	}, nil
}

func (c *Client) BookingGRPCClient() *grpc.Client {
	return c.bookingGrpcClient
}

func (c *Client) Close() error {
	if err := c.bookingGrpcClient.Close(); err != nil {
		return err
	}
	return nil
}

func (c *Client) BookingClient() bookingv1.BookingServiceClient {
	return c.bookingClient
}
