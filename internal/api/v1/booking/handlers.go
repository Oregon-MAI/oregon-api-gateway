package booking

import (
	"io"
	"log/slog"
	"net/http"
	"time"

	bookingclient "github.com/OnYyon/oregon-api-gateway/internal/clients/booking"
	bookingservice "github.com/OnYyon/oregon-api-gateway/internal/services/booking"
	bookingv1 "github.com/Oregon-MAI/oregon-infrastructure/contracts/gen/go/booking"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/protobuf/encoding/protojson"
)

var protoJSONMarshaler = protojson.MarshalOptions{
	EmitUnpopulated: true,
	UseProtoNames:   true,
}

var protoJSONUnmarshaler = protojson.UnmarshalOptions{
	DiscardUnknown: true,
}

type Handler struct {
	service bookingservice.Service
	log     *slog.Logger
	tracer  trace.Tracer
}

func NewHandler(svc *bookingservice.Service, log *slog.Logger) *Handler {
	return &Handler{
		service: *svc,
		log:     log.With(slog.String("component", "booking_handler")),
		tracer:  otel.GetTracerProvider().Tracer("gateway/booking_handler"),
	}
}

func (h *Handler) CreateBooking(c *gin.Context) {
	ctx, span := h.tracer.Start(c.Request.Context(), "BookingHandler.CreateBooking")
	defer span.End()

	c.Request = c.Request.WithContext(ctx)

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": "failed to read body"})
		return
	}

	var req bookingv1.CreateBookingRequest
	if err := protoJSONUnmarshaler.Unmarshal(body, &req); err != nil {
		h.log.Warn("invalid request body format", slog.Any("error", err))
		c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body format"})
		return
	}

	userID := c.GetString("user_id")
	if userID != "" {
		req.UserId = userID
	}

	resp, err := h.service.CreateBooking(ctx, &req)
	if err != nil {
		h.log.Error("failed to create booking", slog.Any("error", err))
		c.JSON(bookingclient.GRPCErrToHTTPStatus(err), map[string]string{"error": err.Error()})
		return
	}

	b, _ := protoJSONMarshaler.Marshal(resp)
	c.Data(http.StatusCreated, "application/json", b)
}

func (h *Handler) GetBooking(c *gin.Context) {
	ctx, span := h.tracer.Start(c.Request.Context(), "BookingHandler.GetBooking")
	defer span.End()

	c.Request = c.Request.WithContext(ctx)

	bookingID := c.Param("booking_id")
	if bookingID == "" {
		c.JSON(http.StatusBadRequest, map[string]string{"error": "booking_id is required"})
		return
	}

	resp, err := h.service.GetBooking(ctx, bookingID)
	if err != nil {
		h.log.Error("failed to get booking", slog.String("booking_id", bookingID), slog.Any("error", err))
		c.JSON(bookingclient.GRPCErrToHTTPStatus(err), map[string]string{"error": err.Error()})
		return
	}

	b, _ := protoJSONMarshaler.Marshal(resp)
	c.Data(http.StatusOK, "application/json", b)
}

func (h *Handler) UserCancelBooking(c *gin.Context) {
	ctx, span := h.tracer.Start(c.Request.Context(), "BookingHandler.UserCancelBooking")
	defer span.End()

	c.Request = c.Request.WithContext(ctx)

	bookingID := c.Param("booking_id")
	if bookingID == "" {
		c.JSON(http.StatusBadRequest, map[string]string{"error": "booking_id is required"})
		return
	}

	resp, err := h.service.UserCancelBooking(ctx, bookingID)
	if err != nil {
		h.log.Error("failed to cancel booking (user)", slog.String("booking_id", bookingID), slog.Any("error", err))
		c.JSON(bookingclient.GRPCErrToHTTPStatus(err), map[string]string{"error": err.Error()})
		return
	}

	b, _ := protoJSONMarshaler.Marshal(resp)
	c.Data(http.StatusOK, "application/json", b)
}

func (h *Handler) AdminCancelBooking(c *gin.Context) {
	ctx, span := h.tracer.Start(c.Request.Context(), "BookingHandler.AdminCancelBooking")
	defer span.End()

	c.Request = c.Request.WithContext(ctx)

	bookingID := c.Param("booking_id")
	if bookingID == "" {
		c.JSON(http.StatusBadRequest, map[string]string{"error": "booking_id is required"})
		return
	}

	resp, err := h.service.AdminCancelBooking(ctx, bookingID)
	if err != nil {
		h.log.Error("failed to cancel booking (admin)", slog.String("booking_id", bookingID), slog.Any("error", err))
		c.JSON(bookingclient.GRPCErrToHTTPStatus(err), map[string]string{"error": err.Error()})
		return
	}

	b, _ := protoJSONMarshaler.Marshal(resp)
	c.Data(http.StatusOK, "application/json", b)
}

func (h *Handler) ListBookingsByUser(c *gin.Context) {
	ctx, span := h.tracer.Start(c.Request.Context(), "BookingHandler.ListBookingsByUser")
	defer span.End()

	c.Request = c.Request.WithContext(ctx)

	userID := c.Query("user_id")
	if userID == "" {
		userID = c.GetString("user_id")
	}

	resp, err := h.service.ListBookingsByUser(ctx, userID)
	if err != nil {
		h.log.Error("failed to list bookings by user", slog.String("user_id", userID), slog.Any("error", err))
		c.JSON(bookingclient.GRPCErrToHTTPStatus(err), map[string]string{"error": err.Error()})
		return
	}

	b, _ := protoJSONMarshaler.Marshal(resp)
	c.Data(http.StatusOK, "application/json", b)
}

func (h *Handler) ListBookingsByResource(c *gin.Context) {
	ctx, span := h.tracer.Start(c.Request.Context(), "BookingHandler.ListBookingsByResource")
	defer span.End()

	c.Request = c.Request.WithContext(ctx)

	resourceID := c.Param("id")
	if resourceID == "" {
		resourceID = c.Param("resource_id")
	}
	if resourceID == "" {
		c.JSON(http.StatusBadRequest, map[string]string{"error": "resource_id is required"})
		return
	}

	var from, to time.Time
	if fromStr := c.Query("from"); fromStr != "" {
		if t, err := time.Parse(time.RFC3339, fromStr); err == nil {
			from = t
		}
	}
	if toStr := c.Query("to"); toStr != "" {
		if t, err := time.Parse(time.RFC3339, toStr); err == nil {
			to = t
		}
	}

	resp, err := h.service.ListBookingsByResource(ctx, resourceID, from, to)
	if err != nil {
		h.log.Error("failed to list bookings by resource", slog.String("resource_id", resourceID), slog.Any("error", err))
		c.JSON(bookingclient.GRPCErrToHTTPStatus(err), map[string]string{"error": err.Error()})
		return
	}

	b, _ := protoJSONMarshaler.Marshal(resp)
	c.Data(http.StatusOK, "application/json", b)
}
