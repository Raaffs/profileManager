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
	health *HealthChecker
}

func connectWithRetry(ctx context.Context, dbURL string) (*pgxpool.Pool, error) {
	var pool *pgxpool.Pool
	var err error

	for i := range 10 {
		log.Printf("Attempting DB connection (attempt %d/10)...", i+1)
		
		pool, err = pgxpool.New(ctx, dbURL)
		if err == nil {
			err = pool.Ping(ctx)
			if err == nil {
				log.Println("PING SUCCESSFUL")
				return pool, nil
			}
		}

		log.Printf("DB not ready: %v. Retrying in 2s...", err)
		if pool != nil {
			pool.Close()
		}
		time.Sleep(2 * time.Second)
	}

	return nil, fmt.Errorf("failed to connect to DB: %w", err)
}

func loadEnv() map[string]string {
    if os.Getenv("DOCKER") != "true" {
        if err := godotenv.Load(".env"); err != nil {
            log.Fatal("No local .env found, skipping")
        }
    }
	log.Println("run time env: ",os.Getenv("DOCKER"))

    envMap := map[string]string{
        env.API_PORT:    os.Getenv(env.API_PORT),
        env.DB_URL:      os.Getenv(env.DB_URL),
        env.CLIENT_PORT: os.Getenv(env.CLIENT_PORT),
        env.JWT_SECRET:  os.Getenv(env.JWT_SECRET),
        env.AES_KEY:     os.Getenv(env.AES_KEY),
    }
    return envMap
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	conn, err := connectWithRetry(ctx,loadEnv()[env.DB_URL]);if err!=nil{
		log.Fatalf("Could not connect to DB: %v", err)
	}

	srv := echo.New()
	app := &Application{
		env:    loadEnv(),
		repo:   store.NewPostgresRepo(conn),
		logger: srv.Logger,
		health: &HealthChecker{status: StatusHealthy},
	}

	app.RegisterRoutes(srv)
	app.LoadMiddleware(srv)
	srv.Logger=app.logger
	
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
		log.Printf("Proper server Shutdown Failed: %+v", err)
	}else{
		log.Println("Server exited")
	}
	
}

