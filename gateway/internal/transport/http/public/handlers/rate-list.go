package handlers

import (
	"io"

	"github.com/gofiber/fiber/v2"
	"github.com/tehrelt/mu/gateway/pkg/pb/ratepb"
	"go.opentelemetry.io/otel/trace"
)

type Rate struct {
	Id          string  `json:"id"`
	Name        string  `json:"name"`
	Rate        float64 `json:"rate"`
	MeasureUnit string  `json:"measureUnit"`
	ServiceType string  `json:"serviceType"`
	CabinetId   string  `json:"cabinetId"`
}

type RateListResponse struct {
	Rates []Rate `json:"rates"`
}

func (rlp *RateListResponse) AddRate(rate Rate) {
	rlp.Rates = append(rlp.Rates, rate)
}

func RateListHandler(rater ratepb.RateServiceClient) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()
		span := trace.SpanFromContext(ctx)

		stream, err := rater.List(ctx, &ratepb.ListRequest{})
		if err != nil {
			return err
		}

		resp := RateListResponse{
			Rates: []Rate{},
		}

		for {
			chunk, err := stream.Recv()
			if err == io.EOF {
				span.AddEvent("EOF")
				break
			}
			if err != nil {
				return err
			}

			span.AddEvent("chunk recv")

			resp.AddRate(Rate{
				Id:          chunk.Id,
				Name:        chunk.Name,
				MeasureUnit: chunk.MeasureUnit,
				Rate:        float64(chunk.Rate) / 100,
				ServiceType: chunk.Type.String(),
			})
		}

		return c.JSON(resp)
	}
}
