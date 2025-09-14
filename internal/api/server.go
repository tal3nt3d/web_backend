package api

import (
	"log"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"web_backend/internal/app/handler"
	"web_backend/internal/app/repository"
)

func StartServer() {
	log.Println("Starting server")

	repo, err := repository.NewRepository()
	if err != nil {
		logrus.Error("Ошибка инициализация репозитория")
	}

	handler := handler.NewHandler(repo)

	r := gin.Default()
	
	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./resources")

	r.GET("/", handler.GetDevices)
	r.GET("/order/:id", handler.GetDevice)
	r.GET("/cart", handler.GetCart)

	r.Run()
	log.Println("Server down")
}