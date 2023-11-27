package main

import (
	"log"
	"os"

	"github.com/bkohler93/jamhubapi/internal/app"
	"github.com/joho/godotenv"
)

func main() {
	if os.Getenv("IN_PROD") == "" {
		err := godotenv.Load()
		if err != nil {
			log.Print("Working without .env file")
			log.Println(err)
		}
	}
	app.RunApp()
}
