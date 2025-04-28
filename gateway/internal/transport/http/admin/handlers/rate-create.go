package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/tehrelt/mu/gateway/pkg/pb/ratepb"
)

type CreateRateRequest struct {
	Name        string  `json:"name"`
	InitialRate float64 `json:"initialRate"`
	MeasureUnit string  `json:"measureUnit"`
}

func RateCreateHandler(rater ratepb.RateServiceClient) fiber.Handler {
	return func(c *fiber.Ctx) error {

		var req CreateRateRequest
		if err := c.BodyParser(&req); err != nil {
			return err
		}

		resp, err := rater.Create(c.UserContext(), &ratepb.CreateRequest{
			Name:        req.Name,
			MeasureUnit: req.MeasureUnit,
			InitialRate: int64(req.InitialRate * 100),
		})
		if err != nil {
			return err
		}

		return c.JSON(resp)
	}
}
