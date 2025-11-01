package app

import (
	"context"
	"net/http"
	"time"
	"log"

	"github.com/Avito-courses/course-go-avito-befragment/internal/routes"
)

type App struct {
	server *http.Server
}

func New(address string) *App {
	return &App{
		server: &http.Server{
			Addr:    address,
			Handler: routes.RegisterRoutes(),
		},
	}
}

func (a *App) Run(ctx context.Context) error {
	go func() {
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("server error: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("Shutting down service-courier") 

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := a.server.Shutdown(shutdownCtx); err != nil {
		return err
	}

	return nil
}

