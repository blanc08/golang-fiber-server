package api

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
)

const (
	authorizationKey        = "authorization"
	authorizationTypeBearer = "bearer"
	AuthorizationPayloadKey = "authorization_payload"
)

func (server *Server) tokenMiddleware() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var accessToken string

		// Try to fetch from header first
		headerAccessToken := ctx.GetReqHeaders()["Authorization"]
		fmt.Println("length header : ", len(headerAccessToken))
		if len(headerAccessToken) > 0 {
			accessToken = headerAccessToken[0]
		}

		// if not exist, try to fetch from cookies
		if len(accessToken) == 0 {
			accessToken = string(ctx.Request().Header.Cookie("Authorization"))
		}

		// Parse bearer token
		fields := strings.Fields(accessToken)
		if len(fields) < 2 {
			return fiber.ErrUnauthorized
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationTypeBearer {
			return fiber.ErrUnauthorized
		}

		tokenString := fields[1]
		payload, err := server.tokenMaker.VerifyToken(tokenString)
		if err != nil {
			return fiber.ErrUnauthorized
		}

		ctx.Locals(AuthorizationPayloadKey, payload)
		return ctx.Next()
	}
}
