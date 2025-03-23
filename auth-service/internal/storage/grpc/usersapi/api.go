package usersapi

import "github.com/tehrelt/moi-uslugi/auth-service/pkg/pb/userspb"

type Api struct {
	client userspb.UserServiceClient
}

func New(client userspb.UserServiceClient) *Api {
	return &Api{client}
}
