package routes

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/OnYyon/oregon-api-gateway/internal/api/v1/booking"
	"github.com/OnYyon/oregon-api-gateway/internal/api/v1/resource"
	bookingclient "github.com/OnYyon/oregon-api-gateway/internal/clients/booking"
	"github.com/OnYyon/oregon-api-gateway/internal/clients/grpc"
	resourceclient "github.com/OnYyon/oregon-api-gateway/internal/clients/resource"
	"github.com/OnYyon/oregon-api-gateway/internal/clients/sso"
	"github.com/OnYyon/oregon-api-gateway/internal/config"
	"github.com/OnYyon/oregon-api-gateway/internal/middlewares"
	bookingservice "github.com/OnYyon/oregon-api-gateway/internal/services/booking"
	resourceservice "github.com/OnYyon/oregon-api-gateway/internal/services/resource"
	"github.com/gin-gonic/gin"
)

func Setup(cfg *config.Config, log *slog.Logger, ssoClient *sso.Client) *http.Server {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	ssoProxy := sso.SSOProxy(cfg.SSO.BaseURL, log)
	resourceClient, err := resourceclient.NewClient(
		grpc.NewConfig(
			grpc.WithTarget(cfg.Resource.PublicTarget),
			grpc.WithTimeout(cfg.Resource.Timeout),
			grpc.WithDialTimeout(cfg.Resource.DialTimeout),
		),
		grpc.NewConfig(
			grpc.WithTarget(cfg.Resource.BookingTarget),
			grpc.WithTimeout(cfg.Resource.Timeout),
			grpc.WithDialTimeout(cfg.Resource.DialTimeout),
		),
		log,
	)
	if err != nil {
		log.Error("failed to create resource client", slog.Any("error", err))
	}
	resourceSvc := resourceservice.NewService(resourceClient)
	resourceHandler := resource.NewHandler(resourceSvc, log)

	bookingClient, err := bookingclient.NewClient(
		grpc.NewConfig(
			grpc.WithTarget(cfg.Booking.Target),
			grpc.WithTimeout(cfg.Booking.Timeout),
			grpc.WithDialTimeout(cfg.Booking.DialTimeout),
		),
		log,
	)
	if err != nil {
		log.Error("failed to create booking client", slog.Any("error", err))
	}
	bookingSvc := bookingservice.NewService(bookingClient)
	bookingHandler := booking.NewHandler(bookingSvc, log)

	r.Use(gin.Recovery())
	r.Use(middlewares.Tracing("api-gateway"))
	r.Use(middlewares.Logging(log))

	pub_auth := r.Group("/api/v1/auth")
	{
		pub_auth.POST("/login", ssoProxy)
		pub_auth.POST("/refresh", ssoProxy)
		pub_auth.POST("/register", ssoProxy)
		pub_auth.POST("/validate", ssoProxy)
	}

	private_auth := r.Group("/api/v1/user")
	private_auth.Use(middlewares.AuthMiddleware(ssoClient, cfg.SSO.JWTSecret, log))
	{
		private_auth.GET("/users", ssoProxy)
		private_auth.GET("/user", ssoProxy)
		private_auth.POST("/change_role", middlewares.RequireRole("admin"), ssoProxy)
		private_auth.DELETE("/delete_user", middlewares.RequireRole("admin"), ssoProxy)
	}

	pub_resource := r.Group("/api/v1/resources")
	pub_resource.Use(middlewares.AuthMiddleware(ssoClient, cfg.SSO.JWTSecret, log))
	{
		pub_resource.GET("", resourceHandler.GetAvailableResources)
		pub_resource.GET("/:id", resourceHandler.GetResource)
		pub_resource.POST("", middlewares.RequireRole("admin"), resourceHandler.CreateResource)
		pub_resource.GET("/list", resourceHandler.GetResourcesList)
		pub_resource.PUT("/:id", middlewares.RequireRole("admin"), resourceHandler.UpdateResource)
		pub_resource.DELETE("/:id", middlewares.RequireRole("admin"), resourceHandler.DeleteResource)
		pub_resource.PATCH("/:id/status", resourceHandler.ChangeResourceStatus)

		pub_resource.GET("/:id/status", resourceHandler.CheckResourceStatus)
		pub_resource.PATCH("/:id/occupancy", resourceHandler.UpdateResourceOccupancy)
		pub_resource.GET("/:id/bookings", func(c *gin.Context) {
			c.Params = append(c.Params, gin.Param{Key: "resource_id", Value: c.Param("id")})
			bookingHandler.ListBookingsByResource(c)
		})
	}

	bookings := r.Group("/api/v1/bookings")
	bookings.Use(middlewares.AuthMiddleware(ssoClient, cfg.SSO.JWTSecret, log))
	{
		bookings.POST("", bookingHandler.CreateBooking)
		bookings.GET("/:booking_id", bookingHandler.GetBooking)
		bookings.POST("/:booking_id/cancel", bookingHandler.UserCancelBooking)
		bookings.GET("", bookingHandler.ListBookingsByUser)
	}

	adminBookings := r.Group("/api/v1/admin/bookings")
	adminBookings.Use(middlewares.AuthMiddleware(ssoClient, cfg.SSO.JWTSecret, log))
	{
		adminBookings.POST("/:booking_id/cancel", middlewares.RequireRole("admin"), bookingHandler.AdminCancelBooking)
	}

	return &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.HTTP.Host, cfg.HTTP.Port),
		Handler:      r,
		ReadTimeout:  cfg.HTTP.ReadTimeout,
		WriteTimeout: cfg.HTTP.WriteTimeout,
		IdleTimeout:  cfg.HTTP.IdleTimeout,
	}
}
