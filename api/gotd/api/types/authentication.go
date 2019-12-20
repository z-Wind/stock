package types

type EASObject struct {
	AccessToken           string `json:"access_token"`
	RefreshToken          string `json:"refresh_token"`
	TokenType             string `json:"token_type"`
	ExpiresIn             int64  `json:"expires_in"`
	Scope                 string `json:"scope"`
	RefreshTokenExpiresIn int64  `json:"refresh_token_expires_in"`
}
