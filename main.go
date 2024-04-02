package main

import (
	"context"
	"fmt"

	"github.com/malav4all/golang-api/application"
)

func main() {
	app := application.New()
	err := app.Start(context.TODO())
	if err != nil {
		fmt.Println("faild to start app:", err)
	}
}
