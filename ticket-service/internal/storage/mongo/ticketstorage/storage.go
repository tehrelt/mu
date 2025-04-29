package ticketstorage

import (
	"context"
	"log/slog"

	"github.com/tehrelt/mu-lib/sl"
	"github.com/tehrelt/mu-lib/tracer"
	"github.com/tehrelt/mu/ticket-service/internal/models"
	"github.com/tehrelt/mu/ticket-service/internal/storage/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	m "go.mongodb.org/mongo-driver/mongo"
	"go.opentelemetry.io/otel"
)

type Storage struct {
	db     *m.Database
	logger *slog.Logger
}

func New(db *m.Database) *Storage {
	return &Storage{
		db:     db,
		logger: slog.With(slog.String("struct", "ticketstorage.Storage")),
	}
}

func (s *Storage) Create(ctx context.Context, ticket models.Ticket) (err error) {
	fn := "ticketstorage.Create"
	log := s.logger.With(sl.Method(fn))

	ctx, span := otel.Tracer(tracer.TracerKey).Start(ctx, fn)
	defer span.End()
	defer func() {
		if err != nil {
			span.RecordError(err)
		}

	}()

	c := s.db.Collection(mongo.TICKETS_COLLECTION)

	marshaled, err := marshalTicket(ticket)
	if err != nil {
		log.Error("failed to marshal ticket", sl.Err(err), slog.Any("ticket", ticket))
		return err
	}
	marshaled.header().ID = new(primitive.ObjectID)
	*marshaled.header().ID = primitive.NewObjectID()

	log.Debug("marshaled ticket to storage model", slog.Any("source", ticket), slog.Any("dst", marshaled))

	res, err := c.InsertOne(ctx, marshaled)
	if err != nil {
		log.Error("failed to insert ticket", sl.Err(err), slog.Any("ticket", ticket))
		return err
	}

	insertedId, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		log.Error("failed to parse insertedId", slog.Any("insertedId", res.InsertedID))
		return err
	}

	log.Info("ticket inserted", slog.String("insertedId", insertedId.Hex()))
	ticket.Header().SetId(insertedId.Hex())

	return nil
}

func (s *Storage) List(ctx context.Context, filters *models.TicketFilters) (<-chan models.Ticket, error) {
	c := s.db.Collection(mongo.TICKETS_COLLECTION)
	fn := "ticketstorage.List"
	log := s.logger.With(slog.String("fn", fn))
	ctx, span := otel.Tracer(tracer.TracerKey).Start(ctx, fn)
	tickets := make(chan models.Ticket)

	cursor, err := c.Find(ctx, marshalFilters(filters))
	if err != nil {
		log.Error("failed to find tickets", sl.Err(err), slog.Any("filters", filters))
		return nil, err
	}

	go func() {
		defer span.End()
		defer close(tickets)
		defer cursor.Close(ctx)

		for cursor.Next(ctx) {
			var header Header
			if err := cursor.Decode(&header); err != nil {
				log.Error("failed to decode header", sl.Err(err))
				span.RecordError(err)
				continue
			}

			log.Debug("decoded header", slog.Any("header", header))

			marshaled, err := factoryTicket(header.Type)
			if err != nil {
				log.Error("failed to create ticket", sl.Err(err))
				span.RecordError(err)
				continue
			}

			if err := cursor.Decode(marshaled); err != nil {
				log.Error("failed to decode ticket", sl.Err(err))
				span.RecordError(err)
				continue
			}

			log.Debug("created ticket", slog.Any("ticket", marshaled))

			ticket, err := unmarshalTicket(marshaled)
			if err != nil {
				log.Error("failed to unmarshal ticket", sl.Err(err))
				span.RecordError(err)
				continue
			}

			tickets <- ticket
		}
	}()

	return tickets, nil
}

func (s *Storage) Find(ctx context.Context, id string) (t models.Ticket, err error) {
	fn := "ticketstorage.Find"
	log := s.logger.With(slog.String("fn", fn))

	ctx, span := otel.Tracer(tracer.TracerKey).Start(ctx, fn)
	defer span.End()

	defer func() {
		if err != nil {
			span.RecordError(err)
		}
	}()

	c := s.db.Collection(mongo.TICKETS_COLLECTION)
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	log.Info("finding ticket", slog.String("_id", id))
	res := c.FindOne(ctx, bson.M{
		"_id": oid,
	})
	if err := res.Err(); err != nil {
		slog.Error("failed to find ticket", sl.Err(err), slog.String("_id", id))
		return nil, err
	}

	var header models.TicketHeader
	if err := res.Decode(&header); err != nil {
		return nil, err
	}

	log.Debug("found ticket header", slog.Any("header", header))

	marshaledTicket, err := factoryTicket(header.TicketType)
	if err != nil {
		return nil, err
	}

	if err := res.Decode(marshaledTicket); err != nil {
		return nil, err
	}

	t, err = unmarshalTicket(marshaledTicket)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (s *Storage) Update(ctx context.Context, id string, newStatus models.TicketStatus) (err error) {
	fn := "ticketstorage.Update"
	log := s.logger.With(slog.String("fn", fn))

	ctx, span := otel.Tracer(tracer.TracerKey).Start(ctx, fn)
	defer span.End()
	defer func() {
		if err != nil {
			span.RecordError(err)
		}
	}()

	c := s.db.Collection(mongo.TICKETS_COLLECTION)

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	log.Debug("updating ticket", slog.Any("new status", newStatus))
	if _, err := c.UpdateOne(
		ctx,
		bson.M{
			"_id": oid,
		},
		bson.M{
			"$set": bson.M{
				"status": newStatus,
			},
		},
	); err != nil {
		slog.Error("failed update ticket", sl.Err(err), slog.Any("ticket_id", id))
		return err
	}

	return nil
}
