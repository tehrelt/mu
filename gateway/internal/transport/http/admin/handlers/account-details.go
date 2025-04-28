package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/tehrelt/mu/gateway/pkg/pb/accountpb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func AccountDetailsHandler(accounter accountpb.AccountServiceClient) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()

		id := c.Params("id")
		if id == "" {
			return fiber.ErrBadRequest
		}

		resp, err := accounter.Find(ctx, &accountpb.FindRequest{
			Id: id,
		})
		if err != nil {
			if status.Code(err) == codes.NotFound {
				return fiber.NewError(404, "account not found")
			}
			return fiber.ErrInternalServerError
		}

		account := Account{
			Id:     resp.Account.Id,
			UserId: resp.Account.UserId,
			House: House{
				Id:      resp.Account.House.Id,
				Address: resp.Account.House.Address,
			},
			Balance:   float64(resp.Account.Balance) / 100,
			CreatedAt: time.Unix(resp.Account.CreatedAt, 0),
		}

		return c.JSON(account)
	}
}
