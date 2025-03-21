package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"

	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

const maxRequestPerSecond = 1000
const maxRequestPerIP = 2

type Handler interface {
	RegisterRoutes(app *fiber.App)
}

type Server struct {
	app    *fiber.App
	port   string
	logger *zap.Logger
}

func New(port string, handler Handler, logger *zap.Logger) Server {
	app := fiber.New(fiber.Config{})
	server := Server{app: app, port: port, logger: logger}

	server.app.Use(recover.New())
	server.app.Use(cors.New())

	server.app.Use(limiter.New(limiter.Config{
		Max:        maxRequestPerSecond,
		Expiration: time.Second,
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Global API rate limit exceeded!",
			})
		},
	}))

	server.app.Use(limiter.New(limiter.Config{
		Max:        maxRequestPerIP,
		Expiration: time.Second,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Too many requests, slow down!",
			})
		},
	}))

	server.addRoutes()

	handler.RegisterRoutes(server.app)

	return server
}

func (s Server) addRoutes() {
	s.app.Get("/health", healthCheck)
	s.app.Get("/metrics", adaptor.HTTPHandler(promhttp.Handler()))
}

func (s Server) Run() {
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		shutdownSignal := <-shutdownChan
		s.logger.Info("Received interrupt signal", zap.String("shutdownSignal", shutdownSignal.String()))

		if err := s.app.Shutdown(); err != nil {
			s.logger.Error("Failed to shutdown gracefully", zap.Error(err))
			return
		}

		s.logger.Info("application shutdown gracefully")
	}()

	err := s.app.Listen(s.port)

	if err != nil {
		s.logger.Panic(err.Error())
	}
}

func (s Server) Stop() {
	if err := s.app.Shutdown(); err != nil {
		s.logger.Error(err.Error())
	}
}

func healthCheck(c *fiber.Ctx) error {
	c.Status(fiber.StatusOK)
	return nil
}
