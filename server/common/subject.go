package common

type AccessTokenSubject struct {
	UserId string `json:"userId"`
	Role   int    `json:"role"`
}

type RefreshTokenSubject struct {
	UserId string `json:"userId"`
}
