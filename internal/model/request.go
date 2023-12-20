package model

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	Password string `json:"password"`
}

type LogoutRequest struct {
	Token string `json:"token"`
}

type RefreshRequest struct {
	Token string `json:"token"`
}
