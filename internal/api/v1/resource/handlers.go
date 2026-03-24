package resource

import (
	"log/slog"
	"net/http"

	clientresource "github.com/OnYyon/oregon-api-gateway/internal/clients/resource"
	"github.com/OnYyon/oregon-api-gateway/internal/services/resource"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type Handler struct {
	service resource.Service
	log     *slog.Logger
	tracer  trace.Tracer
}

func NewHandler(svc *resource.Service, log *slog.Logger) *Handler {
	return &Handler{
		service: *svc,
		log:     log.With(slog.String("component", "resource_handler")),
		tracer:  otel.GetTracerProvider().Tracer("gateway/resource_handler"),
	}
}

func (h *Handler) GetAvailableResources(c *gin.Context) {
	ctx, span := h.tracer.Start(c.Request.Context(), "ResourceHandler.GetAvailableResources")
	defer span.End()

	c.Request = c.Request.WithContext(ctx)

	types := c.QueryArray("type")
	location := c.Query("location")

	resp, err := h.service.GetAvailableResources(ctx, types, location)
	if err != nil {
		h.log.Error("failed to get available resources", slog.Any("error", err))
		c.JSON(clientresource.GRPCErrToHTTPStatus(err), map[string]string{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetResource(c *gin.Context) {
	ctx, span := h.tracer.Start(c.Request.Context(), "ResourceHandler.GetResource")
	defer span.End()

	c.Request = c.Request.WithContext(ctx)

	resourceID := c.Param("id")
	if resourceID == "" {
		c.JSON(http.StatusBadRequest, map[string]string{
			"error": "resource_id_required",
		})
		return
	}

	resp, err := h.service.GetResource(ctx, resourceID)
	if err != nil {
		h.log.Error("failed to get resource",
			slog.String("resource_id", resourceID),
			slog.Any("error", err),
		)
		c.JSON(clientresource.GRPCErrToHTTPStatus(err), map[string]string{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}
