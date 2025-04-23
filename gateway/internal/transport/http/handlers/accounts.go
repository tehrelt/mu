package handlers

import (
	"io"

	"github.com/gofiber/fiber/v2"
	"github.com/tehrelt/mu/gateway/internal/dto"
	"github.com/tehrelt/mu/gateway/internal/transport/http/middlewares"
	"github.com/tehrelt/mu/gateway/pkg/pb/accountpb"
)

type UserAccountsResponse struct {
	Accounts []UserAccount `json:"accounts"`
}

type UserAccount struct {
	Id      string `json:"id"`
	HouseId string `json:"houseId"`
	UserId  string `json:"userId"`
	Balance int64  `json:"balance"`
}

func Accounts(svc accountpb.AccountServiceClient) fiber.Handler {
	return func(c *fiber.Ctx) error {

		profile, ok := c.Locals(middlewares.ProfileLocalKey).(*dto.UserProfile)
		if !ok {
			return c.SendStatus(401)
		}

		stream, err := svc.ListUsersAccounts(c.UserContext(), &accountpb.ListUsersAccountsRequest{
			UserId: profile.Id.String(),
		})
		if err != nil {
			return err
		}

		resp := &UserAccountsResponse{
			Accounts: make([]UserAccount, 0),
		}

		for {
			account, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				return err
			}

			resp.Accounts = append(resp.Accounts, UserAccount{
				Id:      account.Id,
				HouseId: account.HouseId,
				UserId:  account.UserId,
				Balance: account.Balance,
			})
		}

		return c.JSON(resp)
	}
}
