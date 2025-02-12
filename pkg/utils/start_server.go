package utils

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
)

func StartServer(e *echo.Echo) {
	connUrl := fmt.Sprintf("%s:%s", os.Getenv("SERVER_HOST"), os.Getenv("SERVER_PORT"))

	if err := e.Start(connUrl); err != nil && err != http.ErrServerClosed {
		log.Fatalf("shutting down the server: %s", err)
	}
}

func StartServerWithGracefulShutdown(e *echo.Echo) {
	connUrl := fmt.Sprintf("%s:%s", os.Getenv("SERVER_HOST"), os.Getenv("SERVER_PORT"))

	go func() {
		if err := e.Start(connUrl); err != nil && err != http.ErrServerClosed {
			log.Fatalf("shutting down the server: %s", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
}
