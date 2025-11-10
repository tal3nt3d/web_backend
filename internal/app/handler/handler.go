package handler

import (
	"errors"
	"net/http"
	"web_backend/internal/app/repository"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
)

type Handler struct {
	Repository *repository.Repository
}

func NewHandler(r *repository.Repository) *Handler {
	return &Handler{
		Repository: r,
	}
}

// RegisterHandler godoc
// @title Amperage Application API
// @version 1.0
// @description API для управления расчётами нагрузки
// @contact.name API Support
// @contact.url http://localhost:8080
// @contact.email support@amperage.com
// @license.name MIT
// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func (h *Handler) RegisterHandler(router *gin.Engine) {
	api := router.Group("/api/v1")

	unauthorized := api.Group("/")
	unauthorized.POST("/users/signup", h.CreateUser)
	unauthorized.POST("/users/signin", h.SignIn)
	unauthorized.GET("/devices", h.GetDevices)
	unauthorized.GET("/device/:id", h.GetDevice)
	unauthorized.GET("/amperage_application/amperage_application-cart", h.GetAmperageApplicationCart)

	authorized := api.Group("/")
	authorized.Use(h.ModeratorMiddleware(false))
	authorized.POST("/device/create-device", h.CreateDevice)
	authorized.PUT("/device/:id/edit-device", h.EditDevice)
	authorized.DELETE("/device/:id/delete-device", h.DeleteDevice)
	authorized.POST("/device/:id/add-to-amperage_application", h.AddToAmperageApplication)
	authorized.POST("/device/:id/add-photo", h.AddPhoto)

	authorized.GET("/amperage_application/all-amperage_applications", h.GetAllAmperageApplications)
	authorized.GET("/amperage_application/:id", h.GetAmperageApplication)
	authorized.PUT("/amperage_application/:id/edit-amperage_application", h.EditAmperageApplication)
	authorized.PUT("/amperage_application/:id/form-amperage_application", h.FormAmperageApplication)
	authorized.PUT("/amperage_application/:id/finish-amperage_application", h.FinishAmperageApplication)
	authorized.DELETE("/amperage_application/:id/delete-amperage_application", h.DeleteAmperageApplication) 

	authorized.DELETE("/dev_app/:device_id/:amperage_application_id", h.DeleteDeviceFromAmperageApplication)
	authorized.PUT("/dev_app/:device_id/:amperage_application_id", h.EditDeviceFromAmperageApplication)

	authorized.GET("/users/:login/info", h.GetInfo)
	authorized.PUT("/users/:login/info", h.EditInfo)
	authorized.POST("/users/signout", h.SignOut)

	moderator := api.Group("/")
	moderator.Use(h.ModeratorMiddleware(true))
	moderator.PUT("/amperage-calculation/:id/form", h.FormAmperageApplication)

	swaggerURL := ginSwagger.URL("/swagger/doc.json")
	router.Any("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, swaggerURL))
	router.GET("/swagger", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	})
}

func (h *Handler) RegisterStatic(router *gin.Engine) {
	router.LoadHTMLGlob("templates/*")
	router.Static("/styles", "./resources/styles")
	router.Static("/img", "./resources/img")
}

func (h *Handler) errorHandler(ctx *gin.Context, errorStatusCode int, err error) {
	logrus.Error(err.Error())

	var errorMessage string
	switch {
	case errors.Is(err, repository.ErrNotFound):
		errorMessage = "Не найден"
	case errors.Is(err, repository.ErrAlreadyExists):
		errorMessage = "Уже существует"
	case errors.Is(err, repository.ErrNotAllowed):
		errorMessage = "Доступ запрещен"
	case errors.Is(err, repository.ErrNoDraft):
		errorMessage = "Черновик не найден"
	default:
		errorMessage = err.Error()
	}

	ctx.JSON(errorStatusCode, gin.H{
		"status":      "error",
		"description": errorMessage,
	})
}

