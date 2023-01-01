package proto

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Nickname string `json:"nickname"`
}

type LoginResponse struct {
	Ret   int32  `json:"ret"`
	Token string `json:"token"`
}
