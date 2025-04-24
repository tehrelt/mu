package handlers

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/tehrelt/mu-lib/sl"
	"github.com/tehrelt/mu/gateway/pkg/pb/ratepb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func RateDetailsHandler(rater ratepb.RateServiceClient) fiber.Handler {
	return func(c *fiber.Ctx) error {

		id := c.Params("id", "")
		if id == "" {
			return fiber.NewError(fiber.StatusBadRequest, "id is required")
		}

		ctx := c.UserContext()

		resp, err := rater.Find(ctx, &ratepb.FindRequest{Id: id})
		if err != nil {
			slog.Error("failed to find rate", sl.Err(err))
			if e, ok := status.FromError(err); ok {
				if e.Code() == codes.NotFound {
					return fiber.NewError(fiber.StatusNotFound, "rate not found")
				}
			}

			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		rate := &Rate{
			Id:          resp.Id,
			Name:        resp.Name,
			Rate:        float64(resp.Rate) / 100,
			MeasureUnit: resp.MeasureUnit,
		}

		return c.JSON(rate)
	}
}
