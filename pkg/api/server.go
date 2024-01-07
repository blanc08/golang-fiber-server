package api

import (
	"fmt"

	database "github.com/blanc08/stok-gas-management-backend/pkg/database/sqlc"
	"github.com/blanc08/stok-gas-management-backend/pkg/token"
	"github.com/blanc08/stok-gas-management-backend/pkg/util"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Server serves HTTP requests
type Server struct {
	config     util.Config
	pool       *pgxpool.Pool
	store      database.Store
	tokenMaker token.Maker
	app        *fiber.App
	validator  util.XValidator
}

func NewServer(config util.Config, store database.Store, pool *pgxpool.Pool) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.V4SymmetricSecretKeyHex)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:     config,
		pool:       pool,
		store:      store,
		tokenMaker: tokenMaker,
		validator:  *util.NewValidator(),
	}

	server.setupApp()
	return server, nil
}

func (server *Server) setupApp() {
	app := fiber.New(
		fiber.Config{
			ErrorHandler: func(c *fiber.Ctx, err error) error {
				return c.Status(fiber.StatusBadRequest).JSON(util.GlobalErrorHandlerResp{
					Success: false,
					Message: err.Error(),
				})
			},
		},
	)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	v1 := app.Group("/api/v1")
	v1.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
		})
	})

	// authentication
	authRoutes := v1.Group("/auth")
	authRoutes.Post("/login", server.login)
	authRoutes.Post("/register", server.register)

	// authenticated
	authenticatedRoutes := v1.Use(server.tokenMiddleware())

	authenticatedRoutes.Get("/users", func(c *fiber.Ctx) error {
		fmt.Println("authorization payload : ", c.Locals("authorization_payload"))
		return c.SendString("OK")
	})

	server.app = app
}

func (server *Server) Start(address string) error {
	return server.app.Listen(address)
}

func errorResponse(err error) fiber.Error {
	return fiber.Error{
		Code:    500,
		Message: err.Error(),
	}
}
