package sso

import (
	"context"
	"net/http"

	"go.opentelemetry.io/otel/attribute"
)

func (c *Client) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	resp := new(LoginResponse)

	err := c.doRequest(
		ctx,
		http.MethodPost,
		"/api/v1/auth/login",
		req,
		resp,
		"SSO.Login",
		attribute.String("username", req.Username),
	)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
