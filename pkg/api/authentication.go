package api

import (
	"fmt"
	"net/http"
	"time"

	database "github.com/blanc08/stok-gas-management-backend/pkg/database/sqlc"
	"github.com/blanc08/stok-gas-management-backend/pkg/util"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type (
	RegisterRequest struct {
		FirstName string `json:"firstName" validate:"required"`
		LastName  string `json:"lastName" validate:"required"`
		Email     string `json:"email" validate:"required,email"`
		Password  string `json:"password" validate:"required,,min=6"`
	}

	UserResponse struct {
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
		Email     string `json:"email"`
	}

	LoginRequest struct {
		Email    string `json:"email" validate:"required,email"`
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

	if err := server.validator.Validate(request); err != nil {
		return ctx.JSON(fiber.Map{
			"message": "bad request",
			"details": err,
		})
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

	user, err := server.store.CreateUser(ctx.Context(), server.pool, args)
	if err != nil {
		fmt.Println("error while creating user : ", err)
		return fiber.NewError(500, err.Error())
	}

	response := newUserResponse(user)
	return ctx.Status(201).JSON(response)

}

func (server *Server) login(ctx *fiber.Ctx) error {

	var request LoginRequest
	if err := ctx.BodyParser(&request); err != nil {
		return fiber.ErrBadRequest
	}

	if errs := server.validator.Validate(request); len(errs) > 0 {
		return ctx.JSON(fiber.Map{
			"message": "bad request",
			"details": errs,
		})
	}

	user, err := server.store.GetUser(ctx.Context(), server.pool, request.Email)
	if err != nil {
		if err == pgx.ErrNoRows {
			return fiber.NewError(401, "invalid credentials")
		}

		return fiber.NewError(500, err.Error())
	}

	err = util.CheckPassword(request.Password, user.Password)
	if err != nil {
		return fiber.ErrInternalServerError
	}

	accessToken, accessTokenPayload, err := server.tokenMaker.CreateToken(
		user.Email,
		server.config.AccessTokenDuration,
	)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	refreshToken, refreshTokenPayload, err := server.tokenMaker.CreateToken(
		user.Email,
		server.config.RefreshTokenDuration,
	)
	if err != nil {
		return fiber.ErrInternalServerError
	}

	_, err = server.store.CreateSession(ctx.Context(), server.pool, database.CreateSessionParams{
		ID:           refreshTokenPayload.Jti,
		Email:        user.Email,
		RefreshToken: refreshToken,
		UserAgent:    string(ctx.Context().Request.Header.UserAgent()),
		ClientIp:     ctx.Context().RemoteIP().String(),
		IsBlocked:    false,
		ExpiredAt:    refreshTokenPayload.ExpiredAt,
	})
	if err != nil {
		return fiber.ErrInternalServerError
	}

	response := LoginResponse{
		SessionID:             refreshTokenPayload.Jti,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessTokenPayload.ExpiredAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshTokenPayload.ExpiredAt,
		User:                  newUserResponse(user),
	}

	authorizationCookie := new(fiber.Cookie)
	authorizationCookie.Name = "Authorization"
	authorizationCookie.Value = "Bearer " + accessToken
	authorizationCookie.Expires = accessTokenPayload.ExpiredAt

	ctx.Cookie(authorizationCookie)

	return ctx.Status(http.StatusOK).JSON(response)
}
