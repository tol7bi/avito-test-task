package http

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"pvz-backend/internal/models"
	"pvz-backend/internal/repository"
	"pvz-backend/internal/metrics"
)

func CreateReceptionHandler(c *fiber.Ctx) error {
	var req models.CreateReceptionRequest
	if err := c.BodyParser(&req); err != nil || req.PVZID == "" {
		return c.Status(400).JSON(fiber.Map{"message": "невалидный запрос"})
	}

	reception, err := repository.CreateReception(context.Background(), req.PVZID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"message": err.Error()})
	}
	metrics.ReceptionsCreated.Inc()
	return c.Status(201).JSON(reception)
}

func CloseLastReceptionHandler(c *fiber.Ctx) error {
	pvzID := c.Params("pvzId")
	if pvzID == "" {
		return c.Status(400).JSON(fiber.Map{"message": "не указан pvzId"})
	}

	reception, err := repository.CloseLastReception(context.Background(), pvzID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"message": err.Error()})
	}

	return c.Status(200).JSON(reception)
}