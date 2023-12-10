package main

import (
	_ "awesomeProject/docs"
	"awesomeProject/internal/app/config"
	"awesomeProject/internal/app/dsn"
	"awesomeProject/internal/app/handler"
	myminio "awesomeProject/internal/app/myMinio"
	app "awesomeProject/internal/app/pkg"
	"awesomeProject/internal/app/redis"
	"awesomeProject/internal/app/repository"
	"context"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// @title SpaceY
// @version 1.0
// @description Starship's flights

// @contact.name API Support
// @contact.url https://vk.com/bmstu_schedule
// @contact.email bitop@spatecon.ru

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

// @license.name AS IS (NO WARRANTY)

// @host localhost:8080
// @schemes http
// @BasePath /
func main() {
	logger := logrus.New()
	router := gin.Default()
	client := myminio.NewMinioClient(logger)

	conf, err := config.NewConfig()
	if err != nil {
		logger.Fatal( /*"Error of configuration: %s",*/ err)
	}
	ctx := context.Background()
	redis, errRedis := redis.New(ctx, conf.Redis)
	if errRedis != nil {
		logger.Fatalf("Errof with redis connect: %s", err)
	}

	dsn := dsn.FromEnv()
	fmt.Println(dsn)

	repo, err := repository.New(dsn)
	if err != nil {
		logger.Fatalf("Repository error: %s", err)
	}

	handler := handler.New(conf, logger, repo, client, redis)
	application := app.New(conf, router, logger, handler)
	application.Run()
	log.Println("Application start!")
	log.Println("Application terminated!")
}

// $HOME/go/bin/swag
//$HOME/go/bin/swag init -g cmd/main/main.go
