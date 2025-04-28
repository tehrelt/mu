package handlers

import (
	"errors"
	"io"
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/tehrelt/mu-lib/sl"
	"github.com/tehrelt/mu/gateway/pkg/pb/billingpb"
)

type Payment struct {
	Id        string     `json:"id"`
	Status    string     `json:"status"`
	Amount    float64    `json:"amount"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
}

type AccountListPaymentsListResponse struct {
	Payments []Payment `json:"payments"`
}

func (r *AccountListPaymentsListResponse) AddPayment(payment Payment) {
	r.Payments = append(r.Payments, payment)
}

func AccountPaymentsListHandler(biller billingpb.BillingServiceClient) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()

		id := c.Params("id")
		if id == "" {
			return fiber.ErrBadRequest
		}

		stream, err := biller.List(ctx, &billingpb.ListRequest{
			AccountId: id,
		})
		if err != nil {
			slog.Error("failed get list of payments", sl.Err(err))
			return fiber.ErrInternalServerError
		}

		var response AccountListPaymentsListResponse
		response.Payments = make([]Payment, 0, 4)

		for {
			bill, err := stream.Recv()
			if errors.Is(err, io.EOF) {
				break
			}
			if err != nil {
				slog.Error("failed to receive payment", sl.Err(err))
				return fiber.NewError(500, "failed to receive payment")
			}

			payment := Payment{
				Id:        bill.Payment.Id,
				Status:    bill.Payment.Status.String(),
				Amount:    float64(bill.Payment.Amount) / 100,
				CreatedAt: time.Unix(bill.Payment.CreatedAt, 0),
			}

			response.AddPayment(payment)
		}

		return c.JSON(response)
	}
}
