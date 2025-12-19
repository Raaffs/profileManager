package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/Raaffs/profileManager/server/internal/repository"
	"github.com/Raaffs/profileManager/server/internal/store"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

type Application struct {
	env       map[string]string
	repo	 *repository.Repository
}


func main() {

	ctx:= context.Background()
	conn, err := pgxpool.New(ctx, "")
	if err != nil {
		log.Fatalf("Error connecting to database: %v\n", err)
	}
	app:= &Application{
		env: make(map[string]string),
		repo: store.NewPostgresRepo(conn),
	}

	fmt.Println(app)

	srv:=echo.New()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	
	go func() {
		log.Printf("Server is running on %s\n", ":8080")
		if err := srv.Start("8080"); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe(): %v", err)
		}
	}()

	<-stop
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
	log.Println("Server exited properly")
}



