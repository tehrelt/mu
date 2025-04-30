package handlers

import (
	"errors"
	"io"
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/tehrelt/mu-lib/sl"
	"github.com/tehrelt/mu/gateway/pkg/pb/accountpb"
	"github.com/tehrelt/mu/gateway/pkg/pb/consumptionpb"
	"github.com/tehrelt/mu/gateway/pkg/pb/ratepb"
)

type AccountServicesListResponse struct {
	Services []Rate `json:"services"`
}

func AccountServicesListHandler(
	accounter accountpb.AccountServiceClient,
	rateapi ratepb.RateServiceClient,
	cabinetProvider consumptionpb.ConsumptionServiceClient,
) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()

		accountId := c.Params("id")
		if accountId == "" {
			return fiber.ErrBadRequest
		}

		account, err := accounter.Find(ctx, &accountpb.FindRequest{
			Id: accountId,
		})
		if err != nil {
			slog.Error("failed to find account", sl.Err(err))
			return fiber.ErrInternalServerError
		}

		servicesStream, err := rateapi.ListIds(ctx, &ratepb.ListIdsRequest{
			Ids: account.Account.House.ConnectedServiceIds,
		})
		if err != nil {
			slog.Error("failed to list services", sl.Err(err))
			return fiber.ErrInternalServerError
		}

		var response AccountServicesListResponse
		response.Services = make([]Rate, 0, 4)

		for {
			chunk, err := servicesStream.Recv()
			if errors.Is(err, io.EOF) {
				break
			}
			if err != nil {
				slog.Error("failed to receive payment", sl.Err(err))
				return fiber.NewError(500, "failed to receive payment")
			}

			cabinet, err := cabinetProvider.FindCabinet(ctx, &consumptionpb.FindCabinetRequest{
				Criteria: &consumptionpb.FindCabinetRequest_ViaAccount{
					ViaAccount: &consumptionpb.FindCabinetCriteria{
						AccountId: accountId,
						ServiceId: chunk.Id,
					},
				},
			})
			if err != nil {
				slog.Error("failed to find cabinet", sl.Err(err))
				// return fiber.ErrInternalServerError
			}

			rate := Rate{
				Id:          chunk.Id,
				Name:        chunk.Name,
				Rate:        float64(chunk.Rate) / 100,
				MeasureUnit: chunk.MeasureUnit,
				ServiceType: chunk.Type.String(),
				CabinetId:   cabinet.Cabinet.Id,
			}

			response.Services = append(response.Services, rate)
		}

		return c.JSON(response)
	}
}
