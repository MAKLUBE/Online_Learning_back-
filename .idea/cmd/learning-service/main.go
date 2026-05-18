package main

import (
	"log"

	"online-learning-platform/internal/services/learning/app"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}