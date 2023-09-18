package main

import (
	"awesomeProject/internal/app/pkg"
	"log"
)

func main() {
	application = app.New()
	application.Run()
	log.Println("Application start!")
	log.Println("Application terminated!")
}
