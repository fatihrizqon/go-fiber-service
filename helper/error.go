package helper

import (
	"github.com/fatihrizqon/go-fiber-service/internal/presenter/response"
	"github.com/fatihrizqon/go-fiber-service/logger"
	"github.com/gofiber/fiber/v2"
)

func PanicIfError(err error) {
	if err != nil {
		log := logger.GetLogger()
		log.WithError(err).Errorf("Panic occurred: %v", err)
		panic(err)
	}
}

func HandleError(ctx *fiber.Ctx, status int, err error) {
	if err != nil {
		ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{
			Status:  status,
			Message: err.Error(),
			Errors:  err,
		})
	}
}
