package web

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Mahaveer86619/FrameSense/pkg/config"
	"github.com/Mahaveer86619/FrameSense/pkg/handlers"
	"github.com/Mahaveer86619/FrameSense/pkg/services"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type ServerSettings struct {
	Port    string
	Profile string
	e       *echo.Echo
}

func GetServerSettings() ServerSettings {
	return ServerSettings{
		Port:    config.AppConfig.PORT,
		Profile: config.AppConfig.PROFILE,
		e:       NewServer(),
	}
}

func NewServer() *echo.Echo {
	e := echo.New()

	// Standard Echo Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.Secure())
	e.Use(middleware.Gzip())
	e.Use(middleware.CORS())

	// CORS Middleware
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions, http.MethodPatch},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))

	// Services
	healthService := services.NewHealthService()

	// Groups
	api := e.Group("/api/v1")

	// Handlers
	handlers.NewHealthHandler(api, healthService)

	return e
}

func StartServer(e *echo.Echo) {
	settings := GetServerSettings()

	port := settings.Port
	if port == "" {
		port = "7000"
	}

	fmt.Printf("Starting server on port %s in %s mode...\n", port, settings.Profile)

	if err := e.Start(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
