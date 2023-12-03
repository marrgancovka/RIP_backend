package main

import (
	"awesomeProject/internal/app/config"
	"awesomeProject/internal/app/dsn"
	"awesomeProject/internal/app/handler"
	myminio "awesomeProject/internal/app/myMinio"
	app "awesomeProject/internal/app/pkg"
	"awesomeProject/internal/app/repository"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// @title BITOP
// @version 1.0
// @description Bmstu Open IT Platform

// @contact.name API Support
// @contact.url https://vk.com/bmstu_schedule
// @contact.email bitop@spatecon.ru

// @license.name AS IS (NO WARRANTY)

// @host 127.0.0.1
// @schemes https http
// @BasePath /sw
func main() {
	logger := logrus.New()
	router := gin.Default()
	client := myminio.NewMinioClient(logger)
	conf, err := config.NewConfig()
	if err != nil {
		logger.Fatal( /*"Error of configuration: %s",*/ err)
	}

	dsn := dsn.FromEnv()
	fmt.Println(dsn)

	repo, err := repository.New(dsn)
	if err != nil {
		logger.Fatalf("Repository error: %s", err)
	}

	handler := handler.New(conf, logger, repo, client)
	application := app.New(conf, router, logger, handler)
	application.Run()
	log.Println("Application start!")
	log.Println("Application terminated!")
}

// $HOME/go/bin/swag
