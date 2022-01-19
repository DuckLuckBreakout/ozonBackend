package api

type Token struct {
	AccessToken  string `json:"accessToken" valid:"required, type(string)"`
	RefreshToken string `json:"refreshToken" valid:"required, type(string)"`
}
