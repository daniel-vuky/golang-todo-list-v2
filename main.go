package main

import (
	"context"
	"github.com/daniel-vuky/golang-todo-list-and-chat/application"
	"github.com/joho/godotenv"
	"log"
	"os"
	"os/signal"
)

func main() {
	loadEnvErr := godotenv.Load()
	if loadEnvErr != nil {
		log.Fatal("Error loading .env file")
	}
	app := application.New()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	err := app.Start(ctx)
	if err != nil {
		log.Fatalf("")
	}
}
