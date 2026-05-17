package main

import (
	"log"

	"online-learning-platform/internal/services/courses/app"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
