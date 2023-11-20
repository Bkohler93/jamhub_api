package main

import (
	"github.com/bkohler93/jamhubapi/internal/app"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	app.RunApp()
}
