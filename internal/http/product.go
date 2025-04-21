package http

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"pvz-backend/internal/models"
	"pvz-backend/internal/repository"
	"pvz-backend/internal/metrics"
)

func AddProductHandler(c *fiber.Ctx) error {
	var req models.CreateProductRequest
	if err := c.BodyParser(&req); err != nil || req.Type == "" || req.PVZID == "" {
		return c.Status(400).JSON(fiber.Map{"message": "невалидный запрос"})
	}

	product, err := repository.AddProduct(context.Background(), req.PVZID, req.Type)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"message": err.Error()})
	}
	metrics.ProductsCreated.Inc()
	return c.Status(201).JSON(product)
}

func DeleteLastProductHandler(c *fiber.Ctx) error {
	pvzID := c.Params("pvzId")
	if pvzID == "" {
		return c.Status(400).JSON(fiber.Map{"message": "не указан pvzId"})
	}

	err := repository.DeleteLastProduct(context.Background(), pvzID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"message": err.Error()})
	}

	return c.SendStatus(200)
}