package sso

import "encoding/json"

type ValidateRequest struct {
	AccessToken string `json:"access_token"`
}

type ValidateResponse struct {
	Validate bool `json:"is_valid"`
}

func (v *ValidateResponse) UnmarshalJSON(data []byte) error {
	var alias struct {
		IsValid any `json:"is_valid"`
	}
	if err := json.Unmarshal(data, &alias); err != nil {
		return err
	}

	switch val := alias.IsValid.(type) {
	case bool:
		v.Validate = val
	case string:
		v.Validate = (val == "True" || val == "true")
	}
	return nil
}
