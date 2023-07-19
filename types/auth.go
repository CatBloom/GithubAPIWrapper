package types

type AuthReq struct {
	Code string `form:"code" binding:"required"`
}

type AuthRes struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}
