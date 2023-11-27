package main

import (
	"log"

	"github.com/bkohler93/jamhubapi/internal/app"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	app.RunApp()
}
