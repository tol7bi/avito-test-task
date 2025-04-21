package middleware

import (
	"github.com/rs/zerolog"
	"os"
	"time"
	"github.com/gofiber/fiber/v2"
)

var Logger zerolog.Logger

func InitLogger() {
	zerolog.TimeFieldFormat = time.RFC3339

	logFile, err := os.OpenFile("logs/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	output := zerolog.MultiLevelWriter(logFile, os.Stdout) 
	if err == nil {
		Logger = zerolog.New(output).With().Timestamp().Logger()
		
	}

	
}


func FiberLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		duration := time.Since(start)

		Logger.Info().
			Str("method", c.Method()).
			Str("path", c.OriginalURL()).
			Int("status", c.Response().StatusCode()).
			Dur("duration", duration).
			Msg("HTTP request")

		return err
	}
}
