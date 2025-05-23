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
	"github.com/tehrelt/mu/gateway/internal/dto"
	"github.com/tehrelt/mu/gateway/internal/transport/http/admin/handlers"
	"github.com/tehrelt/mu/gateway/internal/transport/http/middlewares"
	"github.com/tehrelt/mu/gateway/pkg/pb/accountpb"
	"github.com/tehrelt/mu/gateway/pkg/pb/authpb"
	"github.com/tehrelt/mu/gateway/pkg/pb/billingpb"
	"github.com/tehrelt/mu/gateway/pkg/pb/ratepb"
	"github.com/tehrelt/mu/gateway/pkg/pb/registerpb"
	"github.com/tehrelt/mu/gateway/pkg/pb/ticketpb"
	"github.com/tehrelt/mu/gateway/pkg/pb/userpb"
)

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type Server struct {
	cfg       *config.Config
	fiber     *fiber.App
	auther    authpb.AuthServiceClient
	register  registerpb.RegisterServiceClient
	accounter accountpb.AccountServiceClient
	rater     ratepb.RateServiceClient
	userapi   userpb.UserServiceClient
	biller    billingpb.BillingServiceClient
	ticketer  ticketpb.TicketServiceClient
}

func New(
	cfg *config.Config,
	auther authpb.AuthServiceClient,
	register registerpb.RegisterServiceClient,
	accounter accountpb.AccountServiceClient,
	rater ratepb.RateServiceClient,
	userapi userpb.UserServiceClient,
	biller billingpb.BillingServiceClient,
	ticketer ticketpb.TicketServiceClient,
) *Server {
	return &Server{
		cfg:       cfg,
		fiber:     fiber.New(),
		auther:    auther,
		register:  register,
		accounter: accounter,
		rater:     rater,
		userapi:   userapi,
		biller:    biller,
		ticketer:  ticketer,
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
					slog.String("path", c.BaseURL()),
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
	s.fiber.Use(middlewares.Trace("admin"))

	token := middlewares.BearerToken()
	authmw := middlewares.Auth(s.auther)(dto.RoleAdmin)

	root := s.fiber.Group("/api")

	auth := root.Group("/auth")
	auth.Put("/refresh", handlers.Refresh(s.auther))
	auth.Post("/login", handlers.Login(s.auther))
	auth.Get("/profile", token, authmw, handlers.Profile(s.auther, s.accounter))
	auth.Post("/logout", token, handlers.Logout(s.auther))

	accounts := root.Group("/accounts")
	accounts.Get("/:id", token, authmw, handlers.AccountDetailsHandler(s.accounter))
	accounts.Get("/:id/payments", token, authmw, handlers.AccountPaymentsListHandler(s.biller))
	accounts.Get("/:id/services", token, authmw, handlers.AccountServicesListHandler(s.accounter, s.rater))

	rates := root.Group("/rates")
	rates.Post("/", token, authmw, handlers.RateCreateHandler(s.rater))
	rates.Get("/", token, authmw, handlers.RateListHandler(s.rater))
	rates.Get("/:id", token, authmw, handlers.RateDetailsHandler(s.rater))

	users := root.Group("/users")
	users.Get("/", token, authmw, handlers.UserListHandler(s.userapi))
	users.Get("/:id", token, authmw, handlers.UserDetailHandler(s.userapi))
	users.Get("/:id/accounts", token, authmw, handlers.UserAccountsList(s.accounter))

	tickets := root.Group("/tickets")
	tickets.Get("/", token, authmw, handlers.TicketListHandler(s.ticketer))
	tickets.Get("/:id", token, authmw, handlers.TicketDetailsHandler(s.ticketer))
	tickets.Patch("/:id", token, authmw, handlers.TicketStatusPatchHandler(s.ticketer))
}

func (s *Server) Run(ctx context.Context) error {

	s.setup()

	host := s.cfg.AdminHttpApi.Host
	port := s.cfg.AdminHttpApi.Port
	addr := fmt.Sprintf("%s:%d", host, port)
	slog.Info("start admin http server", slog.String("addr", addr))

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
