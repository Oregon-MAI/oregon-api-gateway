package sso

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type ValidateRequest struct {
	AccessToken string `json:"access_token"`
}

type RegisterRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Username string `json:"username"`
}

type ValidateResponse struct {
	Validate bool `json:"is_valid"`
}
