package ds

type LoginReq struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type LoginResp struct {
	ExpiresIn   int    `json:"expires_in"`
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

type RegisterReq struct {
	UserName     string
	UserPassword string `json:"pass"`
}

type RegisterResp struct {
	Ok bool `json:"ok"`
}
