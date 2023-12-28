package api

import (
	"fmt"

	database "github.com/blanc08/stok-gas-management-backend/database/sqlc"
	"github.com/blanc08/stok-gas-management-backend/pkg/token"
	"github.com/blanc08/stok-gas-management-backend/pkg/util"
	"github.com/gofiber/fiber/v2"
)

// Server serves HTTP requests
type Server struct {
	config     util.Config
	store      database.Store
	tokenMaker token.Maker
	app        *fiber.App
}

func NewServer(config util.Config, store database.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.V4AsymmetricSecretKeyHex)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	server.setupApp()
	return server, nil
}

func (server *Server) setupApp() {
	// Validator || initialize
	// myValidator := util.NewValidator()

	app := fiber.New(
	// fiber.Config{
	// ErrorHandler: func(c *fiber.Ctx, err error) error {
	// 	return c.Status(fiber.StatusBadRequest).JSON(util.GlobalErrorHandlerResp{
	// 		Success: false,
	// 		Message: err.Error(),
	// 	})
	// },
	// }
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
