package main

import (
	"log"

	"online-learning-platform/internal/services/notifications/app"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
