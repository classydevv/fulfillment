package v1

import (
	"github.com/gofiber/fiber/v2"
)

type responseError struct {
	Error string `json:"error" example:"message"`
}

func errorResponse(ctx *fiber.Ctx, code int, msg string) error {
	return ctx.Status(code).JSON(responseError{msg})
}
