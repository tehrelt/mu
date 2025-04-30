package handlers

import (
	"io"
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/tehrelt/mu-lib/sl"
	"github.com/tehrelt/mu/gateway/internal/dto"
	"github.com/tehrelt/mu/gateway/internal/transport/http/middlewares"
	"github.com/tehrelt/mu/gateway/pkg/pb/accountpb"
)

type UserAccountsResponse struct {
	Accounts []UserAccount `json:"accounts"`
}

type HouseInfo struct {
	Id      string `json:"id"`
	Address string `json:"address"`
}

type UserAccount struct {
	Id      string    `json:"id"`
	UserId  string    `json:"userId"`
	House   HouseInfo `json:"house"`
	Balance float64   `json:"balance"`
}

func Accounts(svc accountpb.AccountServiceClient) fiber.Handler {
	return func(c *fiber.Ctx) error {

		profile, ok := c.Locals(middlewares.ProfileLocalKey).(*dto.UserProfile)
		if !ok {
			return c.SendStatus(401)
		}

		stream, err := svc.List(c.UserContext(), &accountpb.ListRequest{
			UserId: profile.Id.String(),
		})
		if err != nil {
			slog.Error("failed to list users accounts", sl.UUID("userId", profile.Id))
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
				Id:     account.Id,
				UserId: account.UserId,
				House: HouseInfo{
					Id:      account.House.Id,
					Address: account.House.Address,
				},
				Balance: float64(account.Balance) / 100,
			})
		}

		return c.JSON(resp)
	}
}

func Account(svc accountpb.AccountServiceClient) fiber.Handler {
	return func(c *fiber.Ctx) error {
		profile, ok := c.Locals(middlewares.ProfileLocalKey).(*dto.UserProfile)
		if !ok {
			return c.SendStatus(401)
		}

		id := c.Params("id", "")
		if id == "" {
			return c.SendStatus(400)
		}

		response, err := svc.Find(c.UserContext(), &accountpb.FindRequest{
			Id: id,
		})
		if err != nil {
			slog.Error("failed to list users accounts", sl.UUID("userId", profile.Id))
			return err
		}

		account := response.Account
		resp := UserAccount{
			Id:     account.Id,
			UserId: account.UserId,
			House: HouseInfo{
				Id:      account.House.Id,
				Address: account.House.Address,
			},
			Balance: float64(account.Balance) / 100,
		}

		return c.JSON(resp)
	}
}
