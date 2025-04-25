package handlers

import (
	"io"
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/tehrelt/mu/gateway/pkg/pb/accountpb"
)

type House struct {
	Id      string `json:"id"`
	Address string `json:"address"`
}

type Account struct {
	Id        string     `json:"id"`
	House     House      `json:"house"`
	UserId    string     `json:"userId"`
	Balance   float64    `json:"balance"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
}

type UserAccountsListResponse struct {
	Accounts []Account `json:"accounts"`
}

func UserAccountsList(accounter accountpb.AccountServiceClient) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()
		userId := c.Params("id")

		stream, err := accounter.List(ctx, &accountpb.ListRequest{
			UserId: userId,
		})
		if err != nil {
			return err
		}

		resp := UserAccountsListResponse{
			Accounts: make([]Account, 0, 4),
		}

		for {
			account, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				return err
			}

			slog.Debug("house attached to account", slog.String("accId", account.Id), slog.Any("house", account.House))

			acc := Account{
				Id:     account.Id,
				UserId: account.UserId,
				House: House{
					Id:      account.House.Id,
					Address: account.House.Address,
				},
				Balance:   float64(account.Balance) / 100,
				CreatedAt: time.Unix(account.CreatedAt, 0),
			}

			resp.Accounts = append(resp.Accounts, acc)
		}

		return c.JSON(resp)
	}
}
