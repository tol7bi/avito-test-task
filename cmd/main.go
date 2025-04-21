package main

import (
	"log"
	"os"
	"github.com/gofiber/fiber/v2"
	"pvz-backend/internal/models"
	"pvz-backend/internal/app"
	"pvz-backend/internal/grpc"
	"pvz-backend/internal/metrics"
	"pvz-backend/internal/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

func main() {
	middleware.InitLogger()

	go func() {
		metrics.Register()
		http.Handle("/metrics", promhttp.Handler())
		if err := http.ListenAndServe(":9000", nil); err != nil {
			panic(err)
		}
	}()

	if err := models.InitDB(); err != nil {
		log.Fatalf("DB init error: %v", err) 
	}
	go grpc.StartGRPCServer()
	app := fiber.New()
	app.Use(metrics.FiberMiddleware())
	app.Use(middleware.FiberLogger())
	appCfg := app.Group("/")
	app_setup.SetupRoutes(appCfg) 
	

	

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(app.Listen(":" + port))
}
