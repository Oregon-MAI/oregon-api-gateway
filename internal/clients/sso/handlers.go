package sso

import (
	"context"
	"net/http"
)

/*
func (c *Client) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	resp := new(LoginResponse)

	err := c.doRequest(
		ctx, http.MethodPost,
		"/api/v1/auth/login",
		req, resp,
		"SSO.Login",
		attribute.String("username", req.Login),
	)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *Client) Refresh(ctx context.Context, req *RefreshRequest) (*LoginResponse, error) {
	resp := new(LoginResponse)
	err := c.doRequest(
		ctx, http.MethodPost,
		"/api/v1/auth/refresh",
		req,
		resp,
		"SSO.Refresh",
	)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *Client) Register(ctx context.Context, req *RegisterRequest) (*LoginResponse, error) {
	resp := new(LoginResponse)
	err := c.doRequest(
		ctx, http.MethodPost,
		"/api/v1/auth/register",
		req, resp,
		"SSO.Register",
		attribute.String("username", req.Login),
	)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
*/

func (c *Client) Validate(ctx context.Context, req *ValidateRequest) (*ValidateResponse, error) {
	resp := new(ValidateResponse)
	err := c.doRequest(
		ctx, http.MethodPost,
		"/api/v1/auth/validate",
		req, resp,
		"SSO.Validate",
	)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
