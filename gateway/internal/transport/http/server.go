package http

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/tehrelt/mu-lib/sl"
	"github.com/tehrelt/mu/gateway/internal/config"
	"github.com/tehrelt/mu/gateway/internal/transport/http/handlers"
	"github.com/tehrelt/mu/gateway/internal/transport/http/middlewares"
	"github.com/tehrelt/mu/gateway/pkg/pb/authpb"
	"github.com/tehrelt/mu/gateway/pkg/pb/registerpb"
)

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type Server struct {
	cfg      *config.Config
	fiber    *fiber.App
	auther   authpb.AuthServiceClient
	register registerpb.RegisterServiceClient
}

func New(
	cfg *config.Config,
	auther authpb.AuthServiceClient,
	register registerpb.RegisterServiceClient,
) *Server {
	return &Server{
		cfg:      cfg,
		fiber:    fiber.New(),
		auther:   auther,
		register: register,
	}
}

func (s *Server) setup() {
	s.fiber = fiber.New(fiber.Config{
		CaseSensitive: false,
		BodyLimit:     1 << 20,
		AppName:       s.cfg.App.Name,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			resp := ErrorResponse{
				Code:    fiber.StatusInternalServerError,
				Message: err.Error(),
			}

			var e *fiber.Error
			if ok := errors.As(err, &e); ok {

				resp.Code = e.Code
				resp.Message = e.Message
				slog.Error(
					"http error",
					slog.Any("error", resp),
				)
				return c.Status(e.Code).JSON(resp)
			}

			return c.Status(resp.Code).JSON(resp)
		},
	})

	s.fiber.Use(cors.New(cors.Config{
		AllowOrigins:     s.cfg.Cors.AllowedOrigins,
		AllowCredentials: true,
	}))
	s.fiber.Use(logger.New())
	s.fiber.Use(middlewares.Trace)

	token := middlewares.BearerToken()
	authmw := middlewares.Auth(s.auther)

	root := s.fiber.Group("/api")

	auth := root.Group("/auth")
	auth.Post("/register", handlers.Register(s.register))
	auth.Put("/refresh", handlers.Refresh(s.auther))
	auth.Post("/login", handlers.Login(s.auther))
	auth.Get("/profile", middlewares.Cookies(), token, authmw(), handlers.Profile(s.auther))
	auth.Post("/logout", token, handlers.Logout(s.auther))

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
