package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/tehrelt/mu/gateway/pkg/pb/consumptionpb"
)

type Cabinet struct {
	Id        string    `json:"id"`
	AccountId string    `json:"accountId"`
	ServiceId string    `json:"serviceId"`
	Consumed  uint64    `json:"consumed"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func FindCabinet(consumer consumptionpb.ConsumptionServiceClient) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()
		cabinetId := c.Params("cabinetId", "")

		cabinet, err := consumer.FindCabinet(ctx, &consumptionpb.FindCabinetRequest{
			Criteria: &consumptionpb.FindCabinetRequest_Id{
				Id: cabinetId,
			},
		})
		if err != nil {
			return c.SendStatus(fiber.StatusNotFound)
		}

		cab := &Cabinet{
			Id:        cabinet.Cabinet.Id,
			AccountId: cabinet.Cabinet.AccountId,
			ServiceId: cabinet.Cabinet.ServiceId,
			Consumed:  cabinet.Cabinet.Consumed,
			CreatedAt: time.Unix(cabinet.Cabinet.CreatedAt, 0),
			UpdatedAt: time.Unix(cabinet.Cabinet.UpdatedAt, 0),
		}

		return c.JSON(cab)
	}
}
