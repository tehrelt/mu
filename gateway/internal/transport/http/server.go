package http

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/tehrelt/mu-lib/sl"
	"github.com/tehrelt/mu/gateway/internal/config"
	"github.com/tehrelt/mu/gateway/internal/transport/http/handlers"
	"github.com/tehrelt/mu/gateway/internal/transport/http/middlewares"
	"github.com/tehrelt/mu/gateway/pkg/pb/authpb"
)

type Server struct {
	cfg    *config.Config
	fiber  *fiber.App
	auther authpb.AuthServiceClient
}

func New(cfg *config.Config, auther authpb.AuthServiceClient) *Server {
	return &Server{
		cfg:    cfg,
		fiber:  fiber.New(),
		auther: auther,
	}
}

func (s *Server) setup() {
	s.fiber = fiber.New(fiber.Config{
		CaseSensitive: false,
		BodyLimit:     1 << 20,
		AppName:       s.cfg.App.Name,
	})

	s.fiber.Use(logger.New())
	s.fiber.Use(middlewares.Trace)

	root := s.fiber.Group("/api")

	root.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Hello, World",
		})
	})

	auth := root.Group("/auth")
	auth.Post("/login", handlers.Login(s.auther))
	auth.Get("/profile", middlewares.BearerToken(), handlers.Profile(s.auther))
}

func (s *Server) Run(ctx context.Context) error {

	s.setup()

	host := s.cfg.Http.Host
	port := s.cfg.Http.Port
	addr := fmt.Sprintf("%s:%d", host, port)
	slog.Info("start http server", slog.String("addr", addr))

	go func() {
		if err := s.fiber.Listen(addr); err != nil {
			slog.Error(
				"failed to start http server",
				sl.Err(err),
			)
		}
	}()

	<-ctx.Done()
	slog.Info("http server stopped")
	return nil
}
