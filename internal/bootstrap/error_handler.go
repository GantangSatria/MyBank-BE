package bootstrap

import (
	"github.com/gofiber/fiber/v3"
	apperrors "github.com/GantangSatria/MyBank-BE/pkg/errors"
	"github.com/GantangSatria/MyBank-BE/pkg/dto/response"
)

func customErrorHandler(c fiber.Ctx, err error) error {
	if appErr, ok := err.(*apperrors.AppError); ok {
		return c.Status(appErr.HTTPStatus).JSON(response.Error(appErr.Message))
	}
	if fiberErr, ok := err.(*fiber.Error); ok {
		return c.Status(fiberErr.Code).JSON(response.Error(fiberErr.Message))
	}
	return c.Status(fiber.StatusInternalServerError).JSON(response.Error("terjadi kesalahan pada server"))
}