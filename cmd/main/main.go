package main

import (
	"fmt"

	"web_backend/internal/app/config"
	"web_backend/internal/app/dsn"
	"web_backend/internal/app/handler"
	"web_backend/internal/app/repository"
	"web_backend/internal/pkg"
	_ "web_backend/docs"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// @title Amperage Application API
// @version 1.0
// @description API для управления расчётами нагрузки
// @contact.name API Support
// @contact.url http://localhost:8080
// @contact.email support@amperage.com
// @license.name MIT
// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	router := gin.Default()
	conf, err := config.NewConfig()
	if err != nil {
		logrus.Fatalf("error loading config: %v", err)
	}

	postgresString := dsn.FromEnv()
	fmt.Println(postgresString)

	rep, errRep := repository.NewRepository(postgresString)
	if errRep != nil {
		logrus.Fatalf("error initializing repository: %v", errRep)
	}

	hand := handler.NewHandler(rep)

	application := pkg.NewApp(conf, router, hand)
	application.RunApp()
}