package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/time/rate"
)

func (app *Application) LoadMiddleware(e *echo.Echo) {

	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus:  true,
		LogURI:     true,
		LogMethod:  true,
		LogLatency: true,
		LogError:   true,

		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			// Line 1: Basic Request Info
			fmt.Printf("[REQUEST] %s | %s | %s\n",
				v.StartTime.Format("15:04:05"), v.Method, v.URI)

			// Line 2: Performance & Stats
			fmt.Printf("[RESULTS] Status: %d | Latency: %s | IP: %s\n",
				v.Status, v.Latency.String(), v.RemoteIP)

			// Line 3: Errors (only if they exist)
			if v.Error != nil {
				fmt.Printf("[ERROR]   %v\n", v.Error)
			}

			fmt.Println("---")
			return nil
		},
	}))
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowCredentials: true,
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	}))

	config := middleware.RateLimiterConfig{
		Skipper: middleware.DefaultSkipper,
		Store: middleware.NewRateLimiterMemoryStoreWithConfig(
			middleware.RateLimiterMemoryStoreConfig{Rate: rate.Limit(10), Burst: 30, ExpiresIn: 3 * time.Minute},
		),
		IdentifierExtractor: func(ctx echo.Context) (string, error) {
			return ctx.RealIP(), nil
		},
		ErrorHandler: func(context echo.Context, err error) error {
			return context.JSON(http.StatusForbidden, nil)
		},
		DenyHandler: func(context echo.Context, identifier string, err error) error {
			return context.JSON(http.StatusTooManyRequests, nil)
		},
	}
	e.Use(middleware.RateLimiterWithConfig(config))
}

func (app *Application) RegisterRoutes(e *echo.Echo) {
	e.POST("/api/login", app.Login)
	e.POST("/api/register", app.Register)

	e.POST("/api/profile", app.CreateProfile)
	e.GET("/api/profile", app.GetProfile)
	e.PUT("/api/profile", app.UpdateProfile)
}
