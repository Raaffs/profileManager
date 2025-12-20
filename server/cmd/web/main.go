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
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Application struct {
	env    map[string]string
	repo   *repository.Repository
	logger echo.Logger
}

func loadEnv() map[string]string {
	// It reads the .env file and injects it into the OS environment
    err := godotenv.Load(".env.development") 
    if err != nil {
        log.Fatal("Error loading .env file")
    }
	envMap := map[string]string{
		env.API_PORT:    os.Getenv(env.API_PORT),
		env.DB_URL:      os.Getenv(env.DB_URL),
		env.CLIENT_PORT: os.Getenv(env.CLIENT_PORT),
		env.JWT_SECRET:  os.Getenv(env.JWT_SECRET),
		env.AES_KEY:	os.Getenv(env.AES_KEY),
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
	log.Println("app env",app.env)

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

