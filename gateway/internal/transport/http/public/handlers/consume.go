package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/tehrelt/mu/gateway/internal/transport/http/middlewares"
	"github.com/tehrelt/mu/gateway/pkg/pb/accountpb"
	"github.com/tehrelt/mu/gateway/pkg/pb/consumptionpb"
)

type NewConsumeRequest struct {
	Consumed uint64
}

func NewConsume(consumer consumptionpb.ConsumptionServiceClient, accounter accountpb.AccountServiceClient) fiber.Handler {
	return func(c *fiber.Ctx) error {

		user, _ := middlewares.UserFromLocals(c)

		var req NewConsumeRequest
		if err := c.BodyParser(&req); err != nil {
			return err
		}

		cabinetId := c.Params("cabinetId", "")
		if cabinetId == "" {
			return fiber.NewError(fiber.StatusBadRequest, "cabinetId is required")
		}

		if _, err := uuid.Parse(cabinetId); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "cabinetId must be uuid")
		}

		userAccounts, err := listAccounts(c.UserContext(), accounter, user.Id.String())
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		cabinetResponse, err := consumer.FindCabinet(c.UserContext(), &consumptionpb.FindCabinetRequest{
			Criteria: &consumptionpb.FindCabinetRequest_Id{
				Id: cabinetId,
			},
		})

		_, found := lo.Find(userAccounts, func(acc UserAccount) bool {
			return acc.Id == cabinetResponse.Cabinet.AccountId
		})
		if !found {
			return fiber.NewError(fiber.StatusForbidden, "you have no access to this cabinet")
		}

		in := &consumptionpb.ConsumeRequest{
			CabinetId: cabinetId,
			Consumed:  req.Consumed,
		}

		if _, err := consumer.Consume(c.UserContext(), in); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return c.SendStatus(200)
	}
}
