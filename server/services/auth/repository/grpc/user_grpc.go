package userGrpcRepository

import "context"

type UserGrpcRepository interface {
	CreateUserByEmailAndIp(ctx context.Context, email, ip, firstName, lastname string) (*string, error)
}
