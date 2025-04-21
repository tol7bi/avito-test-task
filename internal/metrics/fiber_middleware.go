package metrics

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

func FiberMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		err := c.Next()

		duration := time.Since(start).Seconds()
		method := c.Method()
		path := c.Route().Path 

		HttpRequestsTotal.WithLabelValues(method, path).Inc()
		HttpRequestDuration.WithLabelValues(method, path).Observe(duration)

		return err
	}
}
