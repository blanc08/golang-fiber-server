package api

import (
	"fmt"
	"net/http"
	"time"

	database "github.com/blanc08/stok-gas-management-backend/database/sqlc"
	"github.com/blanc08/stok-gas-management-backend/pkg/util"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type (
	RegisterRequest struct {
		FirstName string `json:"firstName" validate:"required"`
		LastName  string `json:"lastName" validate:"required"`
		Email     string `json:"email" validate:"required"`
		Password  string `json:"password" validate:"required,,min=6"`
	}

	UserResponse struct {
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
		Email     string `json:"email"`
	}

	LoginRequest struct {
		Email    string `json:"email" validate:"required,alphanum"`
		Password string `json:"password" validate:"required,min=6"`
	}

	LoginResponse struct {
		SessionID             uuid.UUID    `json:"session_id"`
		AccessToken           string       `json:"access_token"`
		AccessTokenExpiresAt  time.Time    `json:"access_token_expires_at"`
		RefreshToken          string       `json:"refresh_token"`
		RefreshTokenExpiresAt time.Time    `json:"refresh_token_expires_at"`
		User                  UserResponse `json:"user"`
	}
)

func newUserResponse(user database.User) UserResponse {
	return UserResponse{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
	}
}

func (server *Server) register(ctx *fiber.Ctx) error {
	var request RegisterRequest
	if err := ctx.BodyParser(&request); err != nil {
		fmt.Println("error while parsing request body : ", err)
		return fiber.ErrUnprocessableEntity
	}

	hashdPassword, err := util.HashPassword(request.Password)
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}

	args := database.CreateUserParams{
		FirstName: request.FirstName,
		LastName:  request.LastName,
		Email:     request.Email,
		Password:  hashdPassword,
	}

	user, err := server.store.CreateUser(ctx.Context(), args)
	if err != nil {
		fmt.Println("error while creating user : ", err)
		return fiber.NewError(500, err.Error())
	}

	response := newUserResponse(user)
	return ctx.Status(201).JSON(response)

}

func (server *Server) login(ctx *fiber.Ctx) error {

	var req LoginRequest
	if err := ctx.BodyParser(&req); err != nil {
		return fiber.ErrBadRequest
	}

	user, err := server.store.GetUser(ctx.Context(), req.Email)
	if err != nil {
		if err == pgx.ErrNoRows {
			return fiber.NewError(401, "invalid credentials")
		}

		return fiber.NewError(500, err.Error())
	}

	err = util.CheckPassword(req.Password, user.Password)
	if err != nil {
		// ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		// return
		return fiber.ErrInternalServerError
	}

	accessToken, accessPayload, err := server.tokenMaker.CreateToken(
		user.Email,
		server.config.AccessTokenDuration,
	)
	if err != nil {
		// ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		// return
		return fiber.ErrInternalServerError
	}

	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(
		user.Email,
		server.config.RefreshTokenDuration,
	)
	if err != nil {
		// ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		// return
		return fiber.ErrInternalServerError
	}

	// session, err := server.store.CreateSession(ctx, db.CreateSessionParams{
	// 	ID:           refreshPayload.ID,
	// 	Username:     user.Username,
	// 	RefreshToken: refreshToken,
	// 	UserAgent:    ctx.Request.UserAgent(),
	// 	ClientIp:     ctx.ClientIP(),
	// 	IsBlocked:    false,
	// 	ExpiredAt:    refreshPayload.ExpiredAt,
	// })
	// if err != nil {
	// 	ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	// 	return
	// }

	response := LoginResponse{
		// SessionID:             session.ID,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessPayload.ExpiredAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshPayload.ExpiredAt,
		User:                  newUserResponse(user),
	}

	return ctx.Status(http.StatusOK).JSON(response)
}
