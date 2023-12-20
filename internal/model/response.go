package model

type HealthResponse struct {
	Status string `json:"status"`
}

type BaseResponse struct {
	Success bool   `json:"success,omitempty"`
	Msg     string `json:"msg,omitempty"`
}

type LoginResponse struct {
	BaseResponse
	Token string `json:"token,omitempty"`
}

type RegisterResponse struct {
	BaseResponse
	Token string `json:"token,omitempty"`
}

type RefreshResponse struct {
	BaseResponse
	Token string `json:"token,omitempty"`
}
