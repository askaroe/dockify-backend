package entity

const (
	ClientID     = "115622401"
	ClientSecret = "c1accc4f9244180f19808e3a89e203d061c27d6c187c225ed97359db587c765e"
	RedirectURI  = "https://aokhan.com/huawei/callback"
	State        = "dockify-state-123"
)

type TokenResp struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
}
