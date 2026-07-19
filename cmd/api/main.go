package main

import (
	"context"
	"log"
	"os"
	"smartwaste/internal/database"
	"smartwaste/internal/domain/route"
	"smartwaste/internal/domain/waste"
	"smartwaste/internal/handler"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	dbpool, err := database.NewPostgresPool(ctx)
	if err != nil {
		log.Fatalf("Failed to create database pool: %v", err)
	}
	defer dbpool.Close()
	// Your application logic here
	log.Println("Database connection is succesfuly")

	//
	routeRepo := route.NewPostgresRepository(dbpool)
	routeHandler := handler.NewRouteHandler(routeRepo)
	wasteRepo := waste.NewPostgresRepository(dbpool)
	wasteHandler := handler.NewWasteHandler(wasteRepo)

	app := fiber.New(fiber.Config{
		AppName: "Smart Waste Management API",
	})

	app.Use(cors.New())

	api := app.Group("/api/v1")
	api.Get("/routes/optimal", routeHandler.GetOptimalRoutes)
	api.Post("/wastes", wasteHandler.HandleCreateWaste)
	api.Put("/containers/:id/telemetry", routeHandler.UpdateTelemetry)
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("Server is starting on port %s...", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
