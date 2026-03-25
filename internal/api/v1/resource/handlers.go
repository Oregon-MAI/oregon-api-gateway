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

func (h *Handler) CheckResourceStatus(c *gin.Context) {
	ctx, span := h.tracer.Start(c.Request.Context(), "ResourceHandler.CheckResourceStatus")
	defer span.End()

	c.Request = c.Request.WithContext(ctx)

	resourceID := c.Param("id")
	if resourceID == "" {
		c.JSON(http.StatusBadRequest, map[string]string{
			"error": "resource_id_required",
		})
		return
	}

	resp, err := h.service.CheckResourceStatus(ctx, resourceID)
	if err != nil {
		h.log.Error("failed to check resource status",
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

func (h *Handler) CreateResource(c *gin.Context) {
	ctx, span := h.tracer.Start(c.Request.Context(), "ResourceHandler.CreateResource")
	defer span.End()

	c.Request = c.Request.WithContext(ctx)

	var req clientresource.CreateResourceRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	resp, err := h.service.CreateResource(ctx, &req)
	if err != nil {
		h.log.Error("failed to create resource", slog.Any("error", err))
		c.JSON(clientresource.GRPCErrToHTTPStatus(err), map[string]string{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

func (h *Handler) GetResourcesList(c *gin.Context) {
	ctx, span := h.tracer.Start(c.Request.Context(), "ResourceHandler.GetResourcesList")
	defer span.End()

	c.Request = c.Request.WithContext(ctx)

	types := c.QueryArray("type")

	resp, err := h.service.GetResourcesList(ctx, types)
	if err != nil {
		h.log.Error("failed to get resources list", slog.Any("error", err))
		c.JSON(clientresource.GRPCErrToHTTPStatus(err), map[string]string{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) UpdateResource(c *gin.Context) {
	ctx, span := h.tracer.Start(c.Request.Context(), "ResourceHandler.UpdateResource")
	defer span.End()

	c.Request = c.Request.WithContext(ctx)

	var req clientresource.UpdateResourceRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	resourceID := c.Param("id")
	if resourceID != "" {
		req.ResourceID = resourceID
	}

	resp, err := h.service.UpdateResource(ctx, &req)
	if err != nil {
		h.log.Error("failed to update resource", slog.Any("error", err))
		c.JSON(clientresource.GRPCErrToHTTPStatus(err), map[string]string{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) DeleteResource(c *gin.Context) {
	ctx, span := h.tracer.Start(c.Request.Context(), "ResourceHandler.DeleteResource")
	defer span.End()

	c.Request = c.Request.WithContext(ctx)

	resourceID := c.Param("id")
	if resourceID == "" {
		c.JSON(http.StatusBadRequest, map[string]string{"error": "resource_id_required"})
		return
	}

	success, err := h.service.DeleteResource(ctx, resourceID)
	if err != nil {
		h.log.Error("failed to delete resource", slog.Any("error", err))
		c.JSON(clientresource.GRPCErrToHTTPStatus(err), map[string]string{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, map[string]bool{"success": success})
}

func (h *Handler) ChangeResourceStatus(c *gin.Context) {
	ctx, span := h.tracer.Start(c.Request.Context(), "ResourceHandler.ChangeResourceStatus")
	defer span.End()

	c.Request = c.Request.WithContext(ctx)

	var req clientresource.ChangeResourceStatusRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	resourceID := c.Param("id")
	if resourceID != "" {
		req.ResourceID = resourceID
	}

	resp, err := h.service.ChangeResourceStatus(ctx, &req)
	if err != nil {
		h.log.Error("failed to change resource status", slog.Any("error", err))
		c.JSON(clientresource.GRPCErrToHTTPStatus(err), map[string]string{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) UpdateResourceOccupancy(c *gin.Context) {
	ctx, span := h.tracer.Start(c.Request.Context(), "ResourceHandler.UpdateResourceOccupancy")
	defer span.End()

	c.Request = c.Request.WithContext(ctx)

	var req clientresource.UpdateResourceOccupancyRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	resourceID := c.Param("id")
	if resourceID != "" {
		req.ResourceID = resourceID
	}

	resp, err := h.service.UpdateResourceOccupancy(ctx, &req)
	if err != nil {
		h.log.Error("failed to update resource occupancy", slog.Any("error", err))
		c.JSON(clientresource.GRPCErrToHTTPStatus(err), map[string]string{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}
