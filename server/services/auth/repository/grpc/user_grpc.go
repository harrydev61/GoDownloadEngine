package userGrpcRepository

import "context"

type UserGrpcRepository interface {
	CreateUserByEmailAndIp(ctx context.Context, email, ip string) (*string, error)
}
