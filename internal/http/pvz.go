package http

import (
	"context"
	"time"
	"github.com/gofiber/fiber/v2"
	"pvz-backend/internal/models"
	"pvz-backend/internal/repository"
	"pvz-backend/internal/metrics"
)




func CreatePVZHandler(c *fiber.Ctx) error {
	var req models.CreatePVZRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "невалидный JSON"})
	}

	pvz, err := repository.CreatePVZ(context.Background(), models.DB, req.City)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"message": err.Error()})
	}
	
	metrics.PVZCreated.Inc()
	return c.Status(201).JSON(pvz)
}


func GetPVZListHandler(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)

	var startDate, endDate *time.Time

	if s := c.Query("startDate"); s != "" {
		if parsed, err := time.Parse(time.RFC3339, s); err == nil {
			startDate = &parsed
		}
	}
	if e := c.Query("endDate"); e != "" {
		if parsed, err := time.Parse(time.RFC3339, e); err == nil {
			endDate = &parsed
		}
	}

	data, err := repository.GetPVZList(context.Background(), startDate, endDate, page, limit)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": err.Error()})
	}

	return c.JSON(data)
}
