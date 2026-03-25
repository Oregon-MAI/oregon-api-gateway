package sso

type ValidateRequest struct {
	AccessToken string `json:"access_token"`
}

type ValidateResponse struct {
	Validate bool `json:"is_valid"`
}
