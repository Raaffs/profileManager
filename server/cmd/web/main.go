package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Raaffs/profileManager/server/internal/env"
	"github.com/Raaffs/profileManager/server/internal/repository"
	"github.com/Raaffs/profileManager/server/internal/store/postgres"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Application struct {
	env    map[string]string
	repo   *repository.Repository
	logger echo.Logger
}

func loadEnv() map[string]string {
	envMap := map[string]string{
		env.API_PORT:    os.Getenv(env.API_PORT),
		env.DB_URL:      os.Getenv(env.DB_URL),
		env.CLIENT_PORT: os.Getenv(env.CLIENT_PORT),
		env.JWT_SECRET:  os.Getenv(env.JWT_SECRET),
	}
	return envMap
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	log.Println("Attempting DB connection...")
	conn, err := pgxpool.New(ctx, "postgres://maria:root@localhost:5432/profile_manager")
	if err != nil {
		fmt.Printf("DB ERROR: %v\n", err)
		return
	}

	log.Println("DB connection object created. Testing ping...")
	if err := conn.Ping(ctx); err != nil {
		fmt.Printf("PING ERROR: %v\n", err)
	}
	log.Println("PING SUCCESSFUL")


	srv := echo.New()

		app := &Application{
		env:    loadEnv(),
		repo:   store.NewPostgresRepo(conn),
		logger: srv.Logger,
	}

	app.RegisterRoutes(srv)
	app.LoadMiddleware(srv)
	go func() {
		log.Println(" Server starting on :8080")
		if err := srv.Start(":8080"); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Forcefully shutting down the server: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("\n'Ctrl+C' received, shutting down server...")
	conn.Close()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Proper server Shutdown Failed: %+v", err)
	}
	srv.Logger=app.logger
	log.Println("Server exited")
}

func logger() echo.MiddlewareFunc {
	return  middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
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
	})
}
