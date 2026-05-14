package sso

import (
	"context"
	"net/http"
)

func (c *Client) Validate(ctx context.Context, req *ValidateRequest) (*ValidateResponse, error) {
	resp := new(ValidateResponse)

	headers := map[string]string{
		"Authorization": "Bearer " + req.AccessToken,
	}

	err := c.doRequest(
		ctx, http.MethodPost,
		"/api/v1/auth/validate",
		nil,
		resp,
		headers,
		"SSO.Validate",
	)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
