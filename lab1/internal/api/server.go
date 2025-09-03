package api

import (
	"log"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"lab1/internal/app/handler"
	"lab1/internal/app/repository"
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

	r.GET("/", handler.GetOrders)
	r.GET("/order/:id", handler.GetOrder)
	r.GET("/cart", handler.GetCart)

	r.Run()
	log.Println("Server down")
}