package booking

import (
	"time"
)

type BookingDTO struct {
	BookingID        string    `json:"booking_id"`
	ResourceID       string    `json:"resource_id"`
	UserID           string    `json:"user_id"`
	ResourceName     string    `json:"resource_name"`
	ResourceLocation string    `json:"resource_location"`
	ResourceType     string    `json:"resource_type"`
	StartsAt         time.Time `json:"starts_at"`
	EndsAt           time.Time `json:"ends_at"`
	Status           string    `json:"status"`
	CancelReason     string    `json:"cancel_reason,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type CreateBookingRequestDTO struct {
	ResourceID string    `json:"resource_id"`
	UserID     string    `json:"user_id"`
	StartsAt   time.Time `json:"starts_at"`
	EndsAt     time.Time `json:"ends_at"`
}

type BookingResponseDTO struct {
	Booking BookingDTO `json:"booking"`
}

type ListBookingsResponseDTO struct {
	Bookings []BookingDTO `json:"bookings"`
}
