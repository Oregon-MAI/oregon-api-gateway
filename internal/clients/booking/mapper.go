package booking

import (
	"time"

	bookingv1 "github.com/Oregon-MAI/oregon-infrastructure/contracts/gen/go/booking"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func FromBooking(b *bookingv1.Booking) *BookingDTO {
	if b == nil {
		return nil
	}

	return &BookingDTO{
		BookingID:        b.BookingId,
		ResourceID:       b.ResourceId,
		UserID:           b.UserId,
		ResourceName:     b.ResourceName,
		ResourceLocation: b.ResourceLocation,
		ResourceType:     b.ResourceType,
		StartsAt:         fromTimestamp(b.StartsAt),
		EndsAt:           fromTimestamp(b.EndsAt),
		Status:           b.Status.String(),
		CancelReason:     b.CancelReason,
		CreatedAt:        fromTimestamp(b.CreatedAt),
		UpdatedAt:        fromTimestamp(b.UpdatedAt),
	}
}

func FromCreateBookingResponse(resp *bookingv1.CreateBookingResponse) *BookingResponseDTO {
	if resp == nil || resp.Booking == nil {
		return nil
	}
	return &BookingResponseDTO{
		Booking: *FromBooking(resp.Booking),
	}
}

func FromGetBookingResponse(resp *bookingv1.GetBookingResponse) *BookingResponseDTO {
	if resp == nil || resp.Booking == nil {
		return nil
	}
	return &BookingResponseDTO{
		Booking: *FromBooking(resp.Booking),
	}
}

func FromUserCancelBookingResponse(resp *bookingv1.UserCancelBookingResponse) *BookingResponseDTO {
	if resp == nil || resp.Booking == nil {
		return nil
	}
	return &BookingResponseDTO{
		Booking: *FromBooking(resp.Booking),
	}
}

func FromAdminCancelBookingResponse(resp *bookingv1.AdminCancelBookingResponse) *BookingResponseDTO {
	if resp == nil || resp.Booking == nil {
		return nil
	}
	return &BookingResponseDTO{
		Booking: *FromBooking(resp.Booking),
	}
}

func FromListBookingsByUserResponse(resp *bookingv1.ListBookingsByUserResponse) *ListBookingsResponseDTO {
	if resp == nil {
		return nil
	}
	dto := &ListBookingsResponseDTO{
		Bookings: make([]BookingDTO, 0, len(resp.Bookings)),
	}
	for _, b := range resp.Bookings {
		if mapped := FromBooking(b); mapped != nil {
			dto.Bookings = append(dto.Bookings, *mapped)
		}
	}
	return dto
}

func FromListBookingsByResourceResponse(resp *bookingv1.ListBookingsByResourceResponse) *ListBookingsResponseDTO {
	if resp == nil {
		return nil
	}
	dto := &ListBookingsResponseDTO{
		Bookings: make([]BookingDTO, 0, len(resp.Bookings)),
	}
	for _, b := range resp.Bookings {
		if mapped := FromBooking(b); mapped != nil {
			dto.Bookings = append(dto.Bookings, *mapped)
		}
	}
	return dto
}

func ToCreateBookingRequest(req *CreateBookingRequestDTO) *bookingv1.CreateBookingRequest {
	if req == nil {
		return nil
	}
	return &bookingv1.CreateBookingRequest{
		ResourceId: req.ResourceID,
		UserId:     req.UserID,
		StartsAt:   timestamppb.New(req.StartsAt),
		EndsAt:     timestamppb.New(req.EndsAt),
	}
}

func ToGetBookingRequest(bookingID string) *bookingv1.GetBookingRequest {
	return &bookingv1.GetBookingRequest{
		BookingId: bookingID,
	}
}

func ToUserCancelBookingRequest(bookingID string) *bookingv1.UserCancelBookingRequest {
	return &bookingv1.UserCancelBookingRequest{
		BookingId: bookingID,
	}
}

func ToAdminCancelBookingRequest(bookingID string) *bookingv1.AdminCancelBookingRequest {
	return &bookingv1.AdminCancelBookingRequest{
		BookingId: bookingID,
	}
}

func ToListBookingsByUserRequest(userID string) *bookingv1.ListBookingsByUserRequest {
	return &bookingv1.ListBookingsByUserRequest{
		UserId: userID,
	}
}

func ToListBookingsByResourceRequest(resourceID string, from, to time.Time) *bookingv1.ListBookingsByResourceRequest {
	req := &bookingv1.ListBookingsByResourceRequest{
		ResourceId: resourceID,
	}
	if !from.IsZero() {
		req.From = timestamppb.New(from)
	}
	if !to.IsZero() {
		req.To = timestamppb.New(to)
	}
	return req
}

func fromTimestamp(ts *timestamppb.Timestamp) time.Time {
	if ts == nil {
		return time.Time{}
	}
	return ts.AsTime()
}
