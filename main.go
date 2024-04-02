package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/malav4all/golang-api/application"
)

func main() {
	app := application.New()
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	err := app.Start(ctx)
	if err != nil {
		fmt.Println("faild to start app:", err)
	}
}
