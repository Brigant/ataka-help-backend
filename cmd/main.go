package main

import (
	"log"

	"github.com/baza-trainee/ataka-help-backend/app/rest"
)

func main() {
	if err := rest.SetupAndRun(); err != nil {
		log.Println(err.Error())
	}
}
