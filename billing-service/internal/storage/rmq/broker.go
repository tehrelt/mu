package rmq

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/rabbitmq/amqp091-go"
	"github.com/tehrelt/mu/billing-service/internal/config"
	"github.com/tehrelt/mu/billing-service/internal/dto"
	"github.com/tehrelt/mu/billing-service/rmqmanager"
	"github.com/tehrelt/mu/billing-service/tracer"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

const traceKey = "rmq broker"

type Broker struct {
	cfg     *config.Config
	manager *rmqmanager.RabbitMqManager
}

func New(cfg *config.Config, ch *amqp091.Channel) *Broker {
	return &Broker{
		cfg:     cfg,
		manager: rmqmanager.New(ch),
	}
}

func (b *Broker) PublishStatusChanged(ctx context.Context, event *dto.EventStatusChanged) error {

	event.Timestamp = time.Now()

	ctx, span := otel.Tracer(tracer.TracerKey).Start(ctx, traceKey, trace.WithSpanKind(trace.SpanKindInternal))
	defer span.End()

	j, err := json.Marshal(event)
	if err != nil {
		slog.Error("failed marshal event data", slog.Any("event", event))
		return err
	}

	return b.manager.Publish(ctx, b.cfg.PaymentStatusChanged.Exchange, b.cfg.PaymentStatusChanged.Routing, j)
}
