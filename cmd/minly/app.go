package main

import (
	"log"
	"minly-backend/internal/pkg/app"
)

func main() {
	a, err := app.New()

	if err != nil {
		log.Fatal(err)
	}
	err = a.Run()

	if err != nil {
		log.Fatal(err)
	}
}
