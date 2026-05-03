package booking

import (
	"context"
	"time"

	"github.com/OnYyon/oregon-api-gateway/internal/clients/booking"
	"github.com/OnYyon/oregon-api-gateway/internal/clients/grpc"
	bookingv1 "github.com/Oregon-MAI/oregon-infrastructure/contracts/gen/go/booking"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Service struct {
	client *booking.Client
}

func NewService(client *booking.Client) *Service {
	return &Service{client: client}
}

func (s *Service) Client() *booking.Client {
	return s.client
}

func (s *Service) CreateBooking(ctx context.Context, req *bookingv1.CreateBookingRequest) (*bookingv1.CreateBookingResponse, error) {
	resp, err := grpc.Call(
		ctx,
		s.client.BookingGRPCClient().Conn(),
		s.client.BookingGRPCClient().Log(),
		s.client.BookingGRPCClient().Tracer(),
		s.client.BookingGRPCClient().Timeout(),
		"Booking.CreateBooking",
		func(ctx context.Context, r *bookingv1.CreateBookingRequest) (*bookingv1.CreateBookingResponse, error) {
			return s.client.BookingClient().CreateBooking(ctx, r)
		}, req)

	if err != nil {
		return nil, booking.MapGRPCErr(err)
	}

	return resp, nil
}

func (s *Service) GetBooking(ctx context.Context, bookingID string) (*bookingv1.GetBookingResponse, error) {
	req := &bookingv1.GetBookingRequest{BookingId: bookingID}

	resp, err := grpc.Call(
		ctx,
		s.client.BookingGRPCClient().Conn(),
		s.client.BookingGRPCClient().Log(),
		s.client.BookingGRPCClient().Tracer(),
		s.client.BookingGRPCClient().Timeout(),
		"Booking.GetBooking",
		func(ctx context.Context, r *bookingv1.GetBookingRequest) (*bookingv1.GetBookingResponse, error) {
			return s.client.BookingClient().GetBooking(ctx, r)
		}, req)

	if err != nil {
		return nil, booking.MapGRPCErr(err)
	}

	return resp, nil
}

func (s *Service) UserCancelBooking(ctx context.Context, bookingID string) (*bookingv1.UserCancelBookingResponse, error) {
	req := &bookingv1.UserCancelBookingRequest{BookingId: bookingID}

	resp, err := grpc.Call(
		ctx,
		s.client.BookingGRPCClient().Conn(),
		s.client.BookingGRPCClient().Log(),
		s.client.BookingGRPCClient().Tracer(),
		s.client.BookingGRPCClient().Timeout(),
		"Booking.UserCancelBooking",
		func(ctx context.Context, r *bookingv1.UserCancelBookingRequest) (*bookingv1.UserCancelBookingResponse, error) {
			return s.client.BookingClient().UserCancelBooking(ctx, r)
		}, req)

	if err != nil {
		return nil, booking.MapGRPCErr(err)
	}

	return resp, nil
}

func (s *Service) AdminCancelBooking(ctx context.Context, bookingID string) (*bookingv1.AdminCancelBookingResponse, error) {
	req := &bookingv1.AdminCancelBookingRequest{BookingId: bookingID}

	resp, err := grpc.Call(
		ctx,
		s.client.BookingGRPCClient().Conn(),
		s.client.BookingGRPCClient().Log(),
		s.client.BookingGRPCClient().Tracer(),
		s.client.BookingGRPCClient().Timeout(),
		"Booking.AdminCancelBooking",
		func(ctx context.Context, r *bookingv1.AdminCancelBookingRequest) (*bookingv1.AdminCancelBookingResponse, error) {
			return s.client.BookingClient().AdminCancelBooking(ctx, r)
		}, req)

	if err != nil {
		return nil, booking.MapGRPCErr(err)
	}

	return resp, nil
}

func (s *Service) ListBookingsByUser(ctx context.Context, userID string) (*bookingv1.ListBookingsByUserResponse, error) {
	req := &bookingv1.ListBookingsByUserRequest{UserId: userID}

	resp, err := grpc.Call(
		ctx,
		s.client.BookingGRPCClient().Conn(),
		s.client.BookingGRPCClient().Log(),
		s.client.BookingGRPCClient().Tracer(),
		s.client.BookingGRPCClient().Timeout(),
		"Booking.ListBookingsByUser",
		func(ctx context.Context, r *bookingv1.ListBookingsByUserRequest) (*bookingv1.ListBookingsByUserResponse, error) {
			return s.client.BookingClient().ListBookingsByUser(ctx, r)
		}, req)

	if err != nil {
		return nil, booking.MapGRPCErr(err)
	}

	return resp, nil
}

func (s *Service) ListBookingsByResource(ctx context.Context, resourceID string, from, to time.Time) (*bookingv1.ListBookingsByResourceResponse, error) {
	req := &bookingv1.ListBookingsByResourceRequest{
		ResourceId: resourceID,
	}

	if !from.IsZero() {
		req.From = timestamppb.New(from)
	}
	if !to.IsZero() {
		req.To = timestamppb.New(to)
	}

	resp, err := grpc.Call(
		ctx,
		s.client.BookingGRPCClient().Conn(),
		s.client.BookingGRPCClient().Log(),
		s.client.BookingGRPCClient().Tracer(),
		s.client.BookingGRPCClient().Timeout(),
		"Booking.ListBookingsByResource",
		func(ctx context.Context, r *bookingv1.ListBookingsByResourceRequest) (*bookingv1.ListBookingsByResourceResponse, error) {
			return s.client.BookingClient().ListBookingsByResource(ctx, r)
		}, req)

	if err != nil {
		return nil, booking.MapGRPCErr(err)
	}

	return resp, nil
}
