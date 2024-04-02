package application

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
)

type App struct {
	router http.Handler
	rdb    *redis.Client
}

func New() *App {
	app := &App{
		rdb: redis.NewClient(&redis.Options{
			Addr:     "103.20.214.75:6379",
			Password: "",
			DB:       0,
		}),
	}
	app.loadRoutes()
	return app
}

func (a *App) Start(ctx context.Context) error {
	server := &http.Server{
		Addr:    ":8080",
		Handler: a.router,
	}
	er := a.rdb.Ping(ctx).Err()

	if er != nil {
		return fmt.Errorf("failed to connect to redis: %w", er)
	}

	defer func() {
		if err := a.rdb.Close(); err != nil {
			fmt.Println("failed to close redis", err)
		}
	}()

	fmt.Println("Starting Server...")

	ch := make(chan error, 1)
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			ch <- fmt.Errorf("failed to start server: %w", err)
		}
		close(ch)
	}()

	select {
	case er = <-ch:
		return er
	case <-ctx.Done():
		timeout, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		return server.Shutdown(timeout)
	}

}
