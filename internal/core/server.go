package core

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
)

func StartServer(dbPool *pgxpool.Pool, grpcClient *grpc.ClientConn, port string, appRoutes *chi.Mux) {
	srv := &http.Server{
		Addr:    port,
		Handler: appRoutes,
	}
	log.Printf("Server started on %s\n", srv.Addr)
	serverErr := make(chan error, 1)

	go func() {
		if err := srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
			close(serverErr)
		}
	}()

	waitGracefulShutdown(srv, dbPool, grpcClient, serverErr)

	log.Println("Shutting down...")
}

func waitGracefulShutdown(srv *http.Server, dbPool *pgxpool.Pool, grpcClient *grpc.ClientConn, serverErr <-chan error) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	var reason string

	select {
	case <-ctx.Done():
		reason = "SIGINT or SIGTERM received"
	case err := <-serverErr:
		reason = fmt.Sprintf("Server error: %v", err)
	}

	log.Println("Shutting down...", reason)

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("Error shutting down server: %v", err)
	} else {
		log.Println("Server shutdown")
	}
	dbPool.Close()
	log.Println("Database connection pool closed")

	if grpcClient != nil {
		grpcClient.Close()
		log.Println("GRPC client closed")
	}

	log.Println("Service shutdown")
}
