package main

import (
	"log"

	"korp_backend/internal/modules/billing/bootstrap"
)

func main() {
	app := bootstrap.New()

	log.Printf("billing-service listening on %s", app.Server.Address())
	if err := app.Server.Run(); err != nil {
		log.Fatal(err)
	}
}
