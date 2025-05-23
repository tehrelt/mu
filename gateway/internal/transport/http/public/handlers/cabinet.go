package handlers

import (
	"io"
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/tehrelt/mu/gateway/pkg/pb/consumptionpb"
	"github.com/tehrelt/mu/gateway/pkg/pb/ratepb"
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

func validateUuid(in string) (uuid.UUID, error) {
	id, err := uuid.Parse(in)
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}

type Log struct {
	Id        string    `json:"id"`
	Consumed  uint64    `json:"consumed"`
	CabinetId string    `json:"cabinetId"`
	AccountId string    `json:"accountId"`
	ServiceId string    `json:"serviceId"`
	CreatedAt time.Time `json:"createdAt"`
}

type LogsListResponse struct {
	Logs  []Log  `json:"logs"`
	Total uint64 `json:"total"`
}

func LogsList(consumer consumptionpb.ConsumptionServiceClient) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()
		cabinetId, err := validateUuid(c.Params("cabinetId", ""))
		if err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		limit := c.QueryInt("limit", 10)
		page := c.QueryInt("page", 1)

		stream, err := consumer.Logs(ctx, &consumptionpb.LogsRequest{
			Pagination: &consumptionpb.Pagination{
				Offset: uint64((page - 1) * limit),
				Limit:  uint64(limit),
			},
			CabinetId: cabinetId.String(),
		})
		if err != nil {
			return c.SendStatus(fiber.StatusNotFound)
		}

		metaResp, err := stream.Recv()
		if err != nil {
			return c.SendStatus(fiber.StatusNotFound)
		}

		logs := make([]Log, 0, metaResp.Meta.Total)

		for {
			batch, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				return c.SendStatus(fiber.StatusNotFound)
			}

			for _, log := range batch.Consumptions {
				logs = append(logs, Log{
					Id:        log.Id,
					Consumed:  log.Consumed,
					CabinetId: log.CabinetId,
					AccountId: log.AccountId,
					ServiceId: log.ServiceId,
					CreatedAt: time.Unix(log.CreatedAt, 0),
				})
			}
		}

		return c.JSON(&LogsListResponse{
			Logs:  logs,
			Total: metaResp.Meta.Total,
		})
	}
}

type AggregateLogsList struct {
	Name string `json:"name"`
	id   string
	Logs []Log `json:"logs"`
}

func (a *AggregateLogsList) AddLog(log Log) {
	a.Logs = append(a.Logs, log)
}

type MultipleLogsList struct {
	Items []AggregateLogsList `json:"items"`
}

func LogsListByAccount(consumer consumptionpb.ConsumptionServiceClient, rater ratepb.RateServiceClient) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()
		accountId, err := validateUuid(c.Params("id", ""))
		if err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		limit := c.QueryInt("limit", 10000)
		page := c.QueryInt("page", 1)

		stream, err := consumer.Logs(ctx, &consumptionpb.LogsRequest{
			Pagination: &consumptionpb.Pagination{
				Offset: uint64((page - 1) * limit),
				Limit:  uint64(limit),
			},
			AccountId: accountId.String(),
		})
		if err != nil {
			return c.SendStatus(fiber.StatusNotFound)
		}

		// meta response
		if _, err := stream.Recv(); err != nil {
			return c.SendStatus(fiber.StatusNotFound)
		}

		resp := MultipleLogsList{
			Items: make([]AggregateLogsList, 0),
		}

		for {
			batch, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				return c.SendStatus(fiber.StatusNotFound)
			}

			for _, log := range batch.Consumptions {

				_, idx, found := lo.FindIndexOf(resp.Items, func(item AggregateLogsList) bool {
					return item.id == log.ServiceId
				})

				log := Log{
					Id:        log.Id,
					Consumed:  log.Consumed,
					CabinetId: log.CabinetId,
					AccountId: log.AccountId,
					ServiceId: log.ServiceId,
					CreatedAt: time.Unix(log.CreatedAt, 0),
				}

				if !found {
					service, err := rater.Find(ctx, &ratepb.FindRequest{
						Id: log.ServiceId,
					})
					if err != nil {
						slog.Error("failed to find service, skipping...", slog.String("serviceId", log.ServiceId))
						continue
					}

					resp.Items = append(resp.Items, AggregateLogsList{
						Name: service.Name,
						id:   log.ServiceId,
						Logs: []Log{log},
					})
				} else {
					resp.Items[idx].AddLog(log)
				}
			}
		}

		return c.JSON(resp)
	}
}
