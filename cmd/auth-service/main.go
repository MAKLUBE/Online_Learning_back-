package main

import (
	"log"

	"online-learning-platform/internal/services/auth/app"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
