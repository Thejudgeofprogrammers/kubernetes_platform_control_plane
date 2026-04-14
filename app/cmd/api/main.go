package main

import (
	"context"
	"control_plane/internal/app"
	"control_plane/internal/config"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "control_plane/docs"
)

// @title Control Plane API
// @version 1.0
// @description API для управления API-клиентами в Kubernetes
// @host localhost:8000
// @BasePath /api/v1
func main() {
	env := config.LoadEnv()
	r, rec := app.NewApp(env)

	ctxWorker, cancelWorker := context.WithCancel(context.Background())
	go rec.Run(ctxWorker)


	srv := &http.Server{
		Addr: fmt.Sprintf(":%s", env.Port),
		Handler: r,
	}

	go func() {
		log.Printf("server started on port %s", env.Port)

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen error: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(
		stop,
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	<- stop

	log.Println("shutdown signal received")

	cancelWorker()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server shutdown: %v", err)
	}

	log.Println("server exited properly")
}
