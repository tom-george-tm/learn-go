package main

import (
	"log"

	"urlshortener/internal/config"
	"urlshortener/internal/database"
	"urlshortener/internal/handlers"
	"urlshortener/internal/repository"
	"urlshortener/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("main: config: %v", err)
	}

	db, err := database.Connect(cfg.DBDsn, cfg.IsDevelopment())
	if err != nil {
		log.Fatalf("main: database: %v", err)
	}

	repo := repository.NewURLRepository(db)
	svc := service.NewURLService(repo, cfg.BaseURL)
	hdl := handlers.NewURLHandler(svc)

	app := fiber.New(fiber.Config{
		BodyLimit:    1 * 1024 * 1024,
		ErrorHandler: errorHandler,
	})

	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} ${method} ${path} (${latency})\n",
	}))
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
	}))

	hdl.RegisterRoutes(app)
	app.Get("/health", health)

	addr := ":" + cfg.Port
	log.Printf("main: listening on %s", addr)
	if err := app.Listen(addr); err != nil {
		log.Fatalf("main: server: %v", err)
	}
}

func errorHandler(ctx *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}
	return ctx.Status(code).JSON(fiber.Map{
		"success": false,
		"error":   err.Error(),
	})
}

func health(ctx *fiber.Ctx) error {
	return ctx.JSON(fiber.Map{"status": "ok"})
}
