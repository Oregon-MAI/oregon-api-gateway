package booking

import (
	"errors"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrBookingNotFound    = errors.New("booking not found")
	ErrBookingUnavailable = errors.New("booking service unavailable")
	ErrInvalidArgument    = errors.New("invalid argument")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrPermissionDenied   = errors.New("permission denied")
	ErrDeadlineExceeded   = errors.New("deadline exceeded")
	ErrInternal           = errors.New("internal error")
)

func MapGRPCErr(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, ErrBookingNotFound) ||
		errors.Is(err, ErrBookingUnavailable) ||
		errors.Is(err, ErrInvalidArgument) {
		return err
	}

	st, ok := status.FromError(err)
	if !ok {
		return fmt.Errorf("booking client: %w", err)
	}

	switch st.Code() {
	case codes.NotFound:
		return ErrBookingNotFound
	case codes.Unavailable, codes.DeadlineExceeded:
		return ErrBookingUnavailable
	case codes.InvalidArgument:
		return ErrInvalidArgument
	case codes.Unauthenticated:
		return ErrUnauthorized
	case codes.PermissionDenied:
		return ErrPermissionDenied
	case codes.Internal, codes.Unknown:
		return fmt.Errorf("booking client: %w", ErrInternal)
	default:
		return fmt.Errorf("booking client: %w", err)
	}
}

func GRPCErrToHTTPStatus(err error) int {
	if err == nil {
		return 200
	}

	switch {
	case errors.Is(err, ErrBookingNotFound):
		return 404
	case errors.Is(err, ErrUnauthorized), errors.Is(err, ErrPermissionDenied):
		return 401
	case errors.Is(err, ErrInvalidArgument):
		return 400
	case errors.Is(err, ErrBookingUnavailable), errors.Is(err, ErrDeadlineExceeded):
		return 503
	default:
		return 500
	}
}
