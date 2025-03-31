package usersapi

import "github.com/tehrelt/mu/auth-service/pkg/pb/userpb"

type Api struct {
	client userpb.UserServiceClient
}

func New(client userpb.UserServiceClient) *Api {
	return &Api{client}
}
