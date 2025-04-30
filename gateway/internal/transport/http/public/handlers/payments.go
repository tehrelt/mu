package handlers

import (
	"io"
	"log/slog"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/tehrelt/mu/gateway/internal/dto"
	"github.com/tehrelt/mu/gateway/pkg/pb/billingpb"
)

type PaymentCreateRequest struct {
	AccountId string  `json:"accountId"`
	Amount    float64 `json:"amount"`
}

func PaymentCreateHandler(biller billingpb.BillingServiceClient) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req PaymentCreateRequest
		if err := c.BodyParser(&req); err != nil {
			return err
		}

		resp, err := biller.Create(c.UserContext(), &billingpb.CreateRequest{
			AccountId: req.AccountId,
			Amount:    int64(req.Amount * 100),
		})
		if err != nil {
			return err
		}

		return c.JSON(resp)
	}
}

type PaymentListResponse struct {
	Payments []dto.Payment `json:"payments"`
}

func PaymentListHandler(biller billingpb.BillingServiceClient) fiber.Handler {
	return func(c *fiber.Ctx) error {

		ctx := c.UserContext()

		id := c.Params("id")
		if id == "" {
			return fiber.NewError(400, "invalid id")
		}

		status := c.Query("status", "")
		req := &billingpb.ListRequest{
			AccountId: id,
		}

		limit := c.QueryInt("limit", 50)
		if limit < 0 {
			return fiber.NewError(400, "invalid limit")
		}

		unsignedlimit, _ := strconv.ParseUint(c.Query("limit", ""), 10, 64)

		req.Pagination = &billingpb.Pagination{
			Limit: unsignedlimit,
		}

		if status != "" {
			req.Status = billingpb.PaymentStatus(billingpb.PaymentStatus_value[status])
		}

		slog.Info("request list bills", slog.Any("req", req))

		stream, err := biller.List(ctx, req)
		if err != nil {
			return err
		}

		resp := &PaymentListResponse{
			Payments: make([]dto.Payment, 0, 4),
		}
		for {
			bill, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				return err
			}

			payment := dto.Payment{
				Id:        bill.Payment.Id,
				Status:    bill.Payment.Status.String(),
				Amount:    float64(bill.Payment.Amount) / 100,
				CreatedAt: time.Unix(bill.Payment.CreatedAt, 0),
			}

			resp.Payments = append(resp.Payments, payment)
		}

		return c.JSON(resp)
	}
}

func PaymentFindHandler(biller billingpb.BillingServiceClient) fiber.Handler {
	return func(c *fiber.Ctx) error {

		ctx := c.UserContext()

		id := c.Params("id")
		if id == "" {
			return fiber.NewError(400, "invalid id")
		}

		req := &billingpb.FindRequest{
			Id: id,
		}

		bill, err := biller.Find(ctx, req)
		if err != nil {
			return err
		}

		payment := dto.Payment{
			Id:        bill.Payment.Id,
			Status:    bill.Payment.Status.String(),
			Amount:    float64(bill.Payment.Amount) / 100,
			CreatedAt: time.Unix(bill.Payment.CreatedAt, 0),
		}

		return c.JSON(payment)
	}
}

func PaymentPayHandler(biller billingpb.BillingServiceClient) fiber.Handler {
	return func(c *fiber.Ctx) error {

		ctx := c.UserContext()

		id := c.Params("id")
		if id == "" {
			return fiber.NewError(400, "invalid id")
		}

		req := &billingpb.PayRequest{
			PaymentId: id,
		}

		if _, err := biller.Pay(ctx, req); err != nil {
			return err
		}

		return c.SendStatus(fiber.StatusOK)
	}
}

func PaymentCancelHandler(biller billingpb.BillingServiceClient) fiber.Handler {
	return func(c *fiber.Ctx) error {

		ctx := c.UserContext()

		id := c.Params("id")
		if id == "" {
			return fiber.NewError(400, "invalid id")
		}

		req := &billingpb.CancelRequest{
			PaymentId: id,
		}

		if _, err := biller.Cancel(ctx, req); err != nil {
			return err
		}

		return c.SendStatus(fiber.StatusOK)
	}
}
