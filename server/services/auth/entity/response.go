package entity

type TokenResponse struct {
	AccessToken string `json:"accessToken"`
	// RefreshToken will be used when access token expired
	// to issue new pair access token and refresh token.
	RefreshToken string `json:"refreshToken"`
}
type AuthRes struct {
	UserId string `json:"userId"`
	Email  string `json:"email"`
}
type UserAccessResponse struct {
	Auth  AuthRes       `json:"auth"`
	Token TokenResponse `json:"token"`
}

type RegisterResponse struct {
	UserId   string `json:"userId"`
	Email    string `json:"email"`
	AuthType int    `json:"auth_type"`
}
