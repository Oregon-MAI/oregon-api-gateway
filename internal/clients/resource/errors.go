package resource

import (
	"errors"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrResourceNotFound    = errors.New("resource not found")
	ErrResourceUnavailable = errors.New("resource service unavailable")
	ErrInvalidArgument     = errors.New("invalid argument")
	ErrUnauthorized        = errors.New("unauthorized")
	ErrPermissionDenied    = errors.New("permission denied")
	ErrDeadlineExceeded    = errors.New("deadline exceeded")
	ErrInternal            = errors.New("internal error")
)

func MapGRPCErr(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, ErrResourceNotFound) ||
		errors.Is(err, ErrResourceUnavailable) ||
		errors.Is(err, ErrInvalidArgument) {
		return err
	}

	st, ok := status.FromError(err)
	if !ok {
		return fmt.Errorf("resource client: %w", err)
	}

	switch st.Code() {
	case codes.NotFound:
		return ErrResourceNotFound
	case codes.Unavailable, codes.DeadlineExceeded:
		return ErrResourceUnavailable
	case codes.InvalidArgument:
		return ErrInvalidArgument
	case codes.Unauthenticated:
		return ErrUnauthorized
	case codes.PermissionDenied:
		return ErrPermissionDenied
	case codes.Internal, codes.Unknown:
		return fmt.Errorf("resource client: %w", ErrInternal)
	default:
		return fmt.Errorf("resource client: %w", err)
	}
}

func GRPCErrToHTTPStatus(err error) int {
	if err == nil {
		return 200
	}

	switch {
	case errors.Is(err, ErrResourceNotFound):
		return 404
	case errors.Is(err, ErrUnauthorized), errors.Is(err, ErrPermissionDenied):
		return 401
	case errors.Is(err, ErrInvalidArgument):
		return 400
	case errors.Is(err, ErrResourceUnavailable), errors.Is(err, ErrDeadlineExceeded):
		return 503
	default:
		return 500
	}
}
